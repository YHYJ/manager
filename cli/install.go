/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-14 14:32:16

Description: 子命令 'install' 的实现
*/

package cli

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gookit/color"
	"github.com/pelletier/go-toml"
	"github.com/yhyj/manager/general"
)

// InstallSelfProgram 安装/更新管理程序本身
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
func InstallSelfProgram(configTree *toml.Tree) {
	// 获取配置项
	config, err := general.LoadConfigToStruct(configTree)
	if err != nil {
		color.Error.Println(err)
		return
	}
	// 开始安装
	color.Info.Tips("Install \x1b[3m%s\x1b[0m programs", general.FgCyanText(config.Install.Self.Name))
	color.Info.Tips("%s: %s\n", general.FgWhiteText("Installation path"), general.PrimaryText(config.Install.ProgramPath))
	// 设置代理
	general.SetVariable("http_proxy", config.Variable.HTTPProxy)
	general.SetVariable("https_proxy", config.Variable.HTTPSProxy)

	// 设置进度条参数
	general.ProgressParameters["view"] = "1"
	general.ProgressParameters["sep"] = "-"

	// 使用配置的安装方式进行安装
	switch strings.ToLower(config.Install.Method) {
	case "release":
		// 创建临时目录
		if err := general.CreateDir(config.Install.ReleaseTemp); err != nil {
			color.Error.Println(err)
			return
		}
		// 安装
		name := config.Install.Self.Name
		textLength := 0                                                                                                                                         // 输出文本的长度
		goGithubLatestReleaseTagApi := color.Sprintf(general.GoLatestReleaseTagApiFormat, config.Install.Go.ReleaseApi, config.Install.Go.GithubUsername, name) // 请求远端仓库最新 Tag 的 API
		// 请求 API - GitHub
		body, err := general.RequestApi(goGithubLatestReleaseTagApi)
		if err != nil {
			text := color.Sprintf("%s\n", general.ErrorText(err))
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}
		// 获取远端版本（用于 release 安装方法）
		remoteTag, err := general.GetLatestReleaseTag(body)
		if err != nil {
			text := color.Sprintf("%s\n", general.ErrorText(err))
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}
		// 获取本地版本
		localProgram := filepath.Join(config.Install.ProgramPath, name) // 本地程序路径
		nameArgs := []string{"version", "--only"}                       // 本地程序参数
		localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
		// 比较远端和本地版本
		if remoteTag == localVersion { // 版本一致，则输出无需更新信息
			text := color.Sprintf("%s %s %s %s\n", general.LatestFlag, general.FgGreenText(name), general.FgYellowText(localVersion), general.FgWhiteText(general.LatestVersionMessage))
			color.Printf(text)
			textLength = general.RealLength(text) // 分隔符长度
		} else { // 版本不一致，则安装或更新程序，并输出已安装/更新信息
			// 下载远端文件（如果 Temp 中已有远端文件则删除重新下载）
			goReleaseTempDir := filepath.Join(config.Install.ReleaseTemp, name)
			if general.FileExist(goReleaseTempDir) {
				if err := os.RemoveAll(goReleaseTempDir); err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					return
				}
			}
			// 组装需要的文件的名称
			fileName := general.FileName{}
			// - checksums.txt
			fileName.ChecksumsFile = "checksums.txt"
			// - Archive File
			fileType := func() string {
				if general.Platform == "windows" {
					return "zip"
				}
				return "tar.gz"
			}()
			archiveFileNameWithoutFileType := color.Sprintf("%s_%s_%s_%s", name, remoteTag, general.Platform, general.Arch)
			fileName.ArchiveFile = color.Sprintf("%s.%s", archiveFileNameWithoutFileType, fileType)
			// 获取 Release 文件信息
			filesInfo, err := general.GetReleaseFileInfo(body, fileName)
			if err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
			general.ProgressParameters["action"] = general.DownloadFlag
			general.ProgressParameters["prefix"] = "Download"
			general.ProgressParameters["project"] = color.Sprintf("[%s]", name)
			general.ProgressParameters["fileName"] = color.Sprintf("[%s]", filesInfo.ChecksumsFileInfo.Name)
			general.ProgressParameters["suffix"] = "from Github release:"
			checksumsLocalPath := filepath.Join(config.Install.ReleaseTemp, name, filesInfo.ChecksumsFileInfo.Name) // Checksums 文件本地存储位置
			if err := general.DownloadFile(filesInfo.ChecksumsFileInfo.DownloadUrl, checksumsLocalPath, general.ProgressParameters); err != nil {
				text := color.Error.Sprintf("error -> %s\n", err)
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
			general.ProgressParameters["action"] = general.DownloadFlag
			general.ProgressParameters["prefix"] = "Download"
			general.ProgressParameters["project"] = color.Sprintf("[%s]", name)
			general.ProgressParameters["fileName"] = color.Sprintf("[%s]", filesInfo.ArchiveFileInfo.Name)
			general.ProgressParameters["suffix"] = "from Github release:"
			archiveLocalPath := filepath.Join(config.Install.ReleaseTemp, name, filesInfo.ArchiveFileInfo.Name) // Release 文件本地存储位置
			if err := general.DownloadFile(filesInfo.ArchiveFileInfo.DownloadUrl, archiveLocalPath, general.ProgressParameters); err != nil {
				text := color.Error.Sprintf("error -> %s\n", err)
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
			// 进到下载的远端文件目录
			if err := general.GoToDir(goReleaseTempDir); err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
			// 使用校验文件校验下载的压缩包
			verificationResult, err := general.FileVerification(checksumsLocalPath, archiveLocalPath)
			if err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
			if verificationResult { // 压缩包校验通过
				// 解压压缩包
				err := general.UnzipFile(archiveLocalPath, goReleaseTempDir)
				if err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					return
				}
				archivedProgram := filepath.Join(goReleaseTempDir, archiveFileNameWithoutFileType, name)                // 解压得到的程序
				archivedResourcesFolder := filepath.Join(goReleaseTempDir, archiveFileNameWithoutFileType, "resources") // 解压得到的资源文件夹
				// 检测本地程序是否存在
				if commandErr != nil { // 不存在，安装
					if err := general.InstallFile(archivedProgram, localProgram, 0755); err != nil { // 安装程序
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						return
					} else { // 为已安装的程序设置可执行权限
						if err := os.Chmod(localProgram, 0755); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
					}
					// 安装资源文件 - desktop 文件
					archivedResourcesDesktopFile := filepath.Join(archivedResourcesFolder, "applications", color.Sprintf("%s.desktop", name))   // 解压得到的资源文件 - desktop 文件
					localResourcesDesktopFile := filepath.Join(config.Install.ResourcesPath, "applications", color.Sprintf("%s.desktop", name)) // 本地资源文件 - desktop 文件
					if general.FileExist(archivedResourcesDesktopFile) {
						if err := general.InstallFile(archivedResourcesDesktopFile, localResourcesDesktopFile, 0644); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
					}
					// 安装资源文件 - icon 文件
					archivedResourcesIconFolder := filepath.Join(archivedResourcesFolder, "pixmaps")   // 解压得到的资源文件 - icon 文件夹
					localResourcesIconFolder := filepath.Join(config.Install.ResourcesPath, "pixmaps") // 本地资源文件 - icon 文件夹
					if general.FileExist(archivedResourcesIconFolder) {
						files, err := general.ListFolderFiles(archivedResourcesIconFolder)
						if err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
						if !general.FileExist(localResourcesIconFolder) {
							err := general.CreateDir(localResourcesIconFolder)
							if err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								return
							}
						}
						for _, file := range files {
							archivedResourcesIconFile := filepath.Join(archivedResourcesIconFolder, file) // 解压得到的资源文件 - icon 文件
							localResourcesIconFile := filepath.Join(localResourcesIconFolder, file)       // 本地资源文件 - icon 文件
							if err := general.InstallFile(archivedResourcesIconFile, localResourcesIconFile, 0644); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						}
					}
					// 本次安装结束分隔符
					text := color.Sprintf("%s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(remoteTag), general.FgMagentaText("installed"))
					color.Printf(text)
					textLength = general.RealLength(text) // 分隔符长度
				} else { // 存在，更新
					if err := os.Remove(localProgram); err != nil { // 删除已安装的旧程序
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						return
					}
					if err := general.InstallFile(archivedProgram, localProgram, 0755); err != nil { // 安装新程序
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						return
					} else { // 为已安装的程序设置可执行权限
						if err := os.Chmod(localProgram, 0755); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
					}
					// 安装资源文件 - desktop 文件
					archivedResourcesDesktopFile := filepath.Join(archivedResourcesFolder, "applications", color.Sprintf("%s.desktop", name))   // 解压得到的资源文件 - desktop 文件
					localResourcesDesktopFile := filepath.Join(config.Install.ResourcesPath, "applications", color.Sprintf("%s.desktop", name)) // 本地资源文件 - desktop 文件
					if general.FileExist(archivedResourcesDesktopFile) {
						if err := general.InstallFile(archivedResourcesDesktopFile, localResourcesDesktopFile, 0644); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
					}
					// 安装资源文件 - icon 文件
					archivedResourcesIconFolder := filepath.Join(archivedResourcesFolder, "pixmaps")   // 解压得到的资源文件 - icon 文件夹
					localResourcesIconFolder := filepath.Join(config.Install.ResourcesPath, "pixmaps") // 本地资源文件 - icon 文件夹
					if general.FileExist(archivedResourcesIconFolder) {
						files, err := general.ListFolderFiles(archivedResourcesIconFolder)
						if err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
						if !general.FileExist(localResourcesIconFolder) {
							err := general.CreateDir(localResourcesIconFolder)
							if err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								return
							}
						}
						for _, file := range files {
							archivedResourcesIconFile := filepath.Join(archivedResourcesIconFolder, file) // 解压得到的资源文件 - icon 文件
							localResourcesIconFile := filepath.Join(localResourcesIconFolder, file)       // 本地资源文件 - icon 文件
							if err := general.InstallFile(archivedResourcesIconFile, localResourcesIconFile, 0644); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						}
					}
					// 本次更新结束分隔符
					text := color.Sprintf("%s %s %s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(localVersion), general.FgWhiteText("-->"), general.NoteText(remoteTag), general.FgMagentaText("updated"))
					color.Printf(text)
					textLength = general.RealLength(text) // 分隔符长度
				}
				// 生成/更新自动补全脚本
				for _, completionDir := range config.Install.Go.CompletionDir {
					if general.FileExist(completionDir) {
						generateArgs := []string{"-c", color.Sprintf("%s completion zsh > %s/_%s", localProgram, completionDir, name)}
						if err := general.RunCommand("bash", generateArgs); err != nil {
							text := color.Sprintf("%s %s\n", general.ErrorFlag, general.ErrorText(general.AcsInstallFailedMessage))
							color.Printf(text)
							textLength = general.RealLength(text) // 分隔符长度
							continue
						} else {
							text := color.Sprintf("%s %s\n", general.SuccessFlag, general.SecondaryText(general.AcsInstallSuccessMessage))
							color.Printf(text)
							textLength = general.RealLength(text) // 分隔符长度
							break
						}
					}
				}
			} else { // 压缩包校验失败
				text := color.Error.Sprintf("Archive file verification failed: %s\n", filesInfo.ArchiveFileInfo.Name)
				color.Printf(text)
				textLength = general.RealLength(text) // 分隔符长度
			}
		}
		// 分隔符和延时（延时使输出更加顺畅）
		general.PrintDelimiter(textLength) // 分隔符
		general.Delay(0.1)                 // 0.01s
	case "source":
		// 创建临时目录
		if err := general.CreateDir(config.Install.SourceTemp); err != nil {
			color.Error.Println(err)
			return
		}
		// 安装
		name := config.Install.Self.Name
		textLength := 0                                                                                                                                      // 输出文本的长度
		goGithubLatestSourceTagApi := color.Sprintf(general.GoLatestSourceTagApiFormat, config.Install.Go.GithubApi, config.Install.Go.GithubUsername, name) // 请求远端仓库最新 Tag 的 API
		goGiteaLatestSourceTagApi := color.Sprintf(general.GoLatestSourceTagApiFormat, config.Install.Go.GiteaApi, config.Install.Go.GiteaUsername, name)    // 请求远端仓库最新 Tag 的 API
		// 请求 API - GitHub
		body, err := general.RequestApi(goGithubLatestSourceTagApi)
		if err != nil {
			color.Error.Println(err)
			// 请求 API - Gitea
			body, err = general.RequestApi(goGiteaLatestSourceTagApi)
			if err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
		}
		// 获取远端版本（用于 source 安装方法）
		remoteTag, err := general.GetLatestSourceTag(body)
		if err != nil {
			text := color.Sprintf("%s\n", general.ErrorText(err))
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}
		// 获取本地版本
		localProgram := filepath.Join(config.Install.ProgramPath, name) // 本地程序路径
		nameArgs := []string{"version", "--only"}                       // 本地程序参数
		localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
		// 比较远端和本地版本
		if remoteTag == localVersion { // 版本一致，则输出无需更新信息
			text := color.Sprintf("%s %s %s %s\n", general.LatestFlag, general.FgGreenText(name), general.FgYellowText(localVersion), general.FgWhiteText(general.LatestVersionMessage))
			color.Printf(text)
			textLength = general.RealLength(text) // 分隔符长度
		} else { // 版本不一致，则安装或更新程序，并输出已安装/更新信息
			// 如果 Temp 中已有远端仓库则删除重新克隆
			goSourceTempDir := filepath.Join(config.Install.SourceTemp, name)
			if general.FileExist(goSourceTempDir) {
				if err := os.RemoveAll(goSourceTempDir); err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					return
				}
			}
			// 克隆远端仓库 - GitHub
			goGithubCloneBaseUrl := color.Sprintf("%s/%s", config.Install.Go.GithubUrl, config.Install.Go.GithubUsername) // 远端仓库基础克隆地址（除仓库名）
			color.Printf("%s %s %s %s ", general.DownloadFlag, general.LightText("Clone"), general.FgGreenText(name), "from GitHub")
			if err := general.CloneRepoViaHTTP(config.Install.SourceTemp, goGithubCloneBaseUrl, name); err != nil {
				color.Printf("%s\n", general.ErrorText("error -> ", err))
				// 克隆远端仓库 - Gitea
				goGiteaCloneBaseUrl := color.Sprintf("%s/%s", config.Install.Go.GiteaUrl, config.Install.Go.GiteaUsername) // 远端仓库基础克隆地址（除仓库名）
				color.Printf("%s %s %s %s ", general.DownloadFlag, general.LightText("Clone"), general.FgGreenText(name), "from Gitea")
				if err := general.CloneRepoViaHTTP(config.Install.SourceTemp, goGiteaCloneBaseUrl, name); err != nil {
					text := color.Sprintf("%s\n", general.ErrorText("error -> ", err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					return
				} else {
					color.Println(general.SuccessText("success"))
				}
			} else {
				color.Println(general.SuccessText("success"))
			}
			// 进到下载的远端文件目录
			if err := general.GoToDir(goSourceTempDir); err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
			// 编译生成程序
			if general.FileExist("Makefile") { // Makefile 文件存在则使用 make 编译
				makeArgs := []string{}
				if err := general.RunCommand("make", makeArgs); err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					return
				}
			} else if general.FileExist("main.go") { // Makefile 文件不存在则使用 `go build` 命令编译
				buildArgs := []string{"build", "-trimpath", "-ldflags=-s -w", "-o", name}
				if err := general.RunCommand("go", buildArgs); err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					return
				}
			} else {
				text := color.Error.Sprintf(general.UnableToCompileMessage)
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
			// 检测编译生成的程序是否存在
			compileProgram := filepath.Join(config.Install.SourceTemp, name, config.Install.Go.GeneratePath, name) // 编译生成的程序
			if general.FileExist(compileProgram) {
				// 检测本地程序是否存在
				if commandErr != nil { // 不存在，安装
					if general.FileExist("Makefile") { // Makefile 文件存在则使用 `make install` 命令安装
						makeArgs := []string{"install"}
						if err := general.RunCommand("make", makeArgs); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
					} else { // Makefile 文件不存在则使用自定义函数安装
						if err := general.InstallFile(compileProgram, localProgram, 0755); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						} else {
							// 为已安装的程序设置可执行权限
							if err := os.Chmod(localProgram, 0755); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								return
							}
						}
					}
					// 本次安装结束分隔符
					text := color.Sprintf("%s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(remoteTag), general.FgMagentaText("installed"))
					color.Printf(text)
					textLength = general.RealLength(text) // 分隔符长度
				} else { // 存在，更新
					if general.FileExist("Makefile") { // Makefile 文件存在则使用 `make install` 命令更新
						makeArgs := []string{"install"}
						if err := general.RunCommand("make", makeArgs); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
					} else { // Makefile 文件不存在则使用自定义函数更新
						if err := os.Remove(localProgram); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						}
						if err := general.InstallFile(compileProgram, localProgram, 0755); err != nil {
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							return
						} else {
							// 为已安装的程序设置可执行权限
							if err := os.Chmod(localProgram, 0755); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								return
							}
						}
					}
					// 本次更新结束分隔符
					text := color.Sprintf("%s %s %s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(localVersion), general.FgWhiteText("-->"), general.NoteText(remoteTag), general.FgMagentaText("updated"))
					color.Printf(text)
					textLength = general.RealLength(text) // 分隔符长度
				}
				// 生成/更新自动补全脚本
				for _, completionDir := range config.Install.Go.CompletionDir {
					if general.FileExist(completionDir) {
						generateArgs := []string{"-c", color.Sprintf("%s completion zsh > %s/_%s", localProgram, completionDir, name)}
						if err := general.RunCommand("bash", generateArgs); err != nil {
							text := color.Sprintf("%s %s\n", general.ErrorFlag, general.ErrorText(general.AcsInstallFailedMessage))
							color.Printf(text)
							textLength = general.RealLength(text) // 分隔符长度
							continue
						} else {
							text := color.Sprintf("%s %s\n", general.SuccessFlag, general.SecondaryText(general.AcsInstallSuccessMessage))
							color.Printf(text)
							textLength = general.RealLength(text) // 分隔符长度
							break
						}
					}
				}
			} else {
				text := color.Error.Sprintf("The source file %s does not exist", compileProgram)
				color.Printf(text)
				textLength = general.RealLength(text) // 分隔符长度
			}
		}
		// 分隔符和延时（延时使输出更加顺畅）
		general.PrintDelimiter(textLength) // 分隔符
		general.Delay(0.1)                 // 0.01s
	default:
		text := color.Error.Sprintf("Unsupported installation method '%s': only 'release' and 'source' are supported", config.Install.Method)
		color.Printf(text)
	}
}

// InstallGolangBasedProgram 安装/更新基于 Golang 的程序
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
func InstallGolangBasedProgram(configTree *toml.Tree) {
	// 获取配置项
	config, err := general.LoadConfigToStruct(configTree)
	if err != nil {
		color.Error.Println(err)
		return
	}
	// 开始安装
	color.Info.Tips("Install \x1b[3m%s\x1b[0m programs", general.FgCyanText("golang-based"))
	color.Info.Tips("%s: %s", general.FgWhiteText("Installation path"), general.PrimaryText(config.Install.ProgramPath))
	// 设置代理
	general.SetVariable("http_proxy", config.Variable.HTTPProxy)
	general.SetVariable("https_proxy", config.Variable.HTTPSProxy)

	// 设置进度条参数
	general.ProgressParameters["view"] = "1"
	general.ProgressParameters["sep"] = "-"

	// 使用配置的安装方式进行安装
	switch strings.ToLower(config.Install.Method) {
	case "release":
		// 创建临时目录
		if err := general.CreateDir(config.Install.ReleaseTemp); err != nil {
			color.Error.Println(err)
			return
		}
		// 让用户选择需要安装/更新的程序
		selectedNames, err := general.MultipleSelectionFilter(config.Install.Go.Names)
		if err != nil {
			color.Error.Println(err)
		}
		// 对所选的程序进行排序
		sort.Strings(selectedNames)
		// 遍历所选程序名
		for _, name := range selectedNames {
			textLength := 0                                                                                                                                         // 输出文本的长度
			goGithubLatestReleaseTagApi := color.Sprintf(general.GoLatestReleaseTagApiFormat, config.Install.Go.ReleaseApi, config.Install.Go.GithubUsername, name) // 请求远端仓库最新 Tag 的 API
			// 请求 API - GitHub
			body, err := general.RequestApi(goGithubLatestReleaseTagApi)
			if err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				continue
			}
			// 获取远端版本（用于 release 安装方法）
			remoteTag, err := general.GetLatestReleaseTag(body)
			if err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				continue
			}
			// 获取本地版本
			localProgram := filepath.Join(config.Install.ProgramPath, name) // 本地程序路径
			nameArgs := []string{"version", "--only"}                       // 本地程序参数
			localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
			// 比较远端和本地版本
			if remoteTag == localVersion { // 版本一致，则输出无需更新信息
				text := color.Sprintf("%s %s %s %s\n", general.LatestFlag, general.FgGreenText(name), general.FgYellowText(localVersion), general.FgWhiteText(general.LatestVersionMessage))
				color.Printf(text)
				textLength = general.RealLength(text) // 分隔符长度
			} else { // 版本不一致，则安装或更新程序，并输出已安装/更新信息
				// 下载远端文件（如果 Temp 中已有远端文件则删除重新下载）
				goReleaseTempDir := filepath.Join(config.Install.ReleaseTemp, name)
				if general.FileExist(goReleaseTempDir) {
					if err := os.RemoveAll(goReleaseTempDir); err != nil {
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
				}
				// 组装需要的文件的名称
				fileName := general.FileName{}
				// - checksums.txt
				fileName.ChecksumsFile = "checksums.txt"
				// - Archive File
				fileType := func() string {
					if general.Platform == "windows" {
						return "zip"
					}
					return "tar.gz"
				}()
				archiveFileNameWithoutFileType := color.Sprintf("%s_%s_%s_%s", name, remoteTag, general.Platform, general.Arch)
				fileName.ArchiveFile = color.Sprintf("%s.%s", archiveFileNameWithoutFileType, fileType)
				// 获取 Release 文件信息
				filesInfo, err := general.GetReleaseFileInfo(body, fileName)
				if err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
				general.ProgressParameters["action"] = general.DownloadFlag
				general.ProgressParameters["prefix"] = "Download"
				general.ProgressParameters["project"] = color.Sprintf("[%s]", name)
				general.ProgressParameters["fileName"] = color.Sprintf("[%s]", filesInfo.ChecksumsFileInfo.Name)
				general.ProgressParameters["suffix"] = "from Github release:"
				checksumsLocalPath := filepath.Join(config.Install.ReleaseTemp, name, filesInfo.ChecksumsFileInfo.Name) // Checksums 文件本地存储位置
				if err := general.DownloadFile(filesInfo.ChecksumsFileInfo.DownloadUrl, checksumsLocalPath, general.ProgressParameters); err != nil {
					text := color.Error.Sprintf("error -> %s\n", err)
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
				general.ProgressParameters["action"] = general.DownloadFlag
				general.ProgressParameters["prefix"] = "Download"
				general.ProgressParameters["project"] = color.Sprintf("[%s]", name)
				general.ProgressParameters["fileName"] = color.Sprintf("[%s]", filesInfo.ArchiveFileInfo.Name)
				general.ProgressParameters["suffix"] = "from Github release:"
				archiveLocalPath := filepath.Join(config.Install.ReleaseTemp, name, filesInfo.ArchiveFileInfo.Name) // Release 文件本地存储位置
				if err := general.DownloadFile(filesInfo.ArchiveFileInfo.DownloadUrl, archiveLocalPath, general.ProgressParameters); err != nil {
					text := color.Error.Sprintf("error -> %s\n", err)
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
				// 进到下载的远端文件目录
				if err := general.GoToDir(goReleaseTempDir); err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
				// 使用校验文件校验下载的压缩包
				verificationResult, err := general.FileVerification(checksumsLocalPath, archiveLocalPath)
				if err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
				if verificationResult { // 压缩包校验通过
					// 解压压缩包
					err := general.UnzipFile(archiveLocalPath, goReleaseTempDir)
					if err != nil {
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
					archivedProgram := filepath.Join(goReleaseTempDir, archiveFileNameWithoutFileType, name)                // 解压得到的程序
					archivedResourcesFolder := filepath.Join(goReleaseTempDir, archiveFileNameWithoutFileType, "resources") // 解压得到的资源文件夹
					// 检测本地程序是否存在
					if commandErr != nil { // 不存在，安装
						if err := general.InstallFile(archivedProgram, localProgram, 0755); err != nil { // 安装程序
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							continue
						} else { // 为已安装的程序设置可执行权限
							if err := os.Chmod(localProgram, 0755); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						}
						// 安装资源文件 - desktop 文件
						archivedResourcesDesktopFile := filepath.Join(archivedResourcesFolder, "applications", color.Sprintf("%s.desktop", name))   // 解压得到的资源文件 - desktop 文件
						localResourcesDesktopFile := filepath.Join(config.Install.ResourcesPath, "applications", color.Sprintf("%s.desktop", name)) // 本地资源文件 - desktop 文件
						if general.FileExist(archivedResourcesDesktopFile) {
							if err := general.InstallFile(archivedResourcesDesktopFile, localResourcesDesktopFile, 0644); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						}
						// 安装资源文件 - icon 文件
						archivedResourcesIconFolder := filepath.Join(archivedResourcesFolder, "pixmaps")   // 解压得到的资源文件 - icon 文件夹
						localResourcesIconFolder := filepath.Join(config.Install.ResourcesPath, "pixmaps") // 本地资源文件 - icon 文件夹
						if general.FileExist(archivedResourcesIconFolder) {
							files, err := general.ListFolderFiles(archivedResourcesIconFolder)
							if err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
							if !general.FileExist(localResourcesIconFolder) {
								err := general.CreateDir(localResourcesIconFolder)
								if err != nil {
									text := color.Sprintf("%s\n", general.ErrorText(err))
									color.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								}
							}
							for _, file := range files {
								archivedResourcesIconFile := filepath.Join(archivedResourcesIconFolder, file) // 解压得到的资源文件 - icon 文件
								localResourcesIconFile := filepath.Join(localResourcesIconFolder, file)       // 本地资源文件 - icon 文件
								if err := general.InstallFile(archivedResourcesIconFile, localResourcesIconFile, 0644); err != nil {
									text := color.Sprintf("%s\n", general.ErrorText(err))
									color.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								}
							}
						}
						// 本次安装结束分隔符
						text := color.Sprintf("%s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(remoteTag), general.FgMagentaText("installed"))
						color.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					} else { // 存在，更新
						if err := os.Remove(localProgram); err != nil { // 删除已安装的旧程序
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							continue
						}
						if err := general.InstallFile(archivedProgram, localProgram, 0755); err != nil { // 安装新程序
							text := color.Sprintf("%s\n", general.ErrorText(err))
							color.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							continue
						} else { // 为已安装的程序设置可执行权限
							if err := os.Chmod(localProgram, 0755); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						}
						// 安装资源文件 - desktop 文件
						archivedResourcesDesktopFile := filepath.Join(archivedResourcesFolder, "applications", color.Sprintf("%s.desktop", name))   // 解压得到的资源文件 - desktop 文件
						localResourcesDesktopFile := filepath.Join(config.Install.ResourcesPath, "applications", color.Sprintf("%s.desktop", name)) // 本地资源文件 - desktop 文件
						if general.FileExist(archivedResourcesDesktopFile) {
							if err := general.InstallFile(archivedResourcesDesktopFile, localResourcesDesktopFile, 0644); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						}
						// 安装资源文件 - icon 文件
						archivedResourcesIconFolder := filepath.Join(archivedResourcesFolder, "pixmaps")   // 解压得到的资源文件 - icon 文件夹
						localResourcesIconFolder := filepath.Join(config.Install.ResourcesPath, "pixmaps") // 本地资源文件 - icon 文件夹
						if general.FileExist(archivedResourcesIconFolder) {
							files, err := general.ListFolderFiles(archivedResourcesIconFolder)
							if err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
							if !general.FileExist(localResourcesIconFolder) {
								err := general.CreateDir(localResourcesIconFolder)
								if err != nil {
									text := color.Sprintf("%s\n", general.ErrorText(err))
									color.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								}
							}
							for _, file := range files {
								archivedResourcesIconFile := filepath.Join(archivedResourcesIconFolder, file) // 解压得到的资源文件 - icon 文件
								localResourcesIconFile := filepath.Join(localResourcesIconFolder, file)       // 本地资源文件 - icon 文件
								if err := general.InstallFile(archivedResourcesIconFile, localResourcesIconFile, 0644); err != nil {
									text := color.Sprintf("%s\n", general.ErrorText(err))
									color.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								}
							}
						}
						// 本次更新结束分隔符
						text := color.Sprintf("%s %s %s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(localVersion), general.FgWhiteText("-->"), general.NoteText(remoteTag), general.FgMagentaText("updated"))
						color.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					}
					// 生成/更新自动补全脚本
					for _, completionDir := range config.Install.Go.CompletionDir {
						if general.FileExist(completionDir) {
							generateArgs := []string{"-c", color.Sprintf("%s completion zsh > %s/_%s", localProgram, completionDir, name)}
							if err := general.RunCommand("bash", generateArgs); err != nil {
								text := color.Sprintf("%s %s\n", general.ErrorFlag, general.ErrorText(general.AcsInstallFailedMessage))
								color.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
								continue
							} else {
								text := color.Sprintf("%s %s\n", general.SuccessFlag, general.SecondaryText(general.AcsInstallSuccessMessage))
								color.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
								break
							}
						}
					}
				} else { // 压缩包校验失败
					text := color.Error.Sprintf("Archive file verification failed: %s\n", filesInfo.ArchiveFileInfo.Name)
					color.Printf(text)
					textLength = general.RealLength(text) // 分隔符长度
				}
			}
			// 分隔符和延时（延时使输出更加顺畅）
			general.PrintDelimiter(textLength) // 分隔符
			general.Delay(0.1)                 // 0.01s
		}
	case "source":
		// 创建临时目录
		if err := general.CreateDir(config.Install.SourceTemp); err != nil {
			color.Error.Println(err)
			return
		}
		// 让用户选择需要安装/更新的程序
		selectedNames, err := general.MultipleSelectionFilter(config.Install.Go.Names)
		if err != nil {
			color.Error.Println(err)
		}
		// 对所选的程序进行排序
		sort.Strings(selectedNames)
		// 遍历所选程序名
		for _, name := range selectedNames {
			textLength := 0                                                                                                                                      // 输出文本的长度
			goGithubLatestSourceTagApi := color.Sprintf(general.GoLatestSourceTagApiFormat, config.Install.Go.GithubApi, config.Install.Go.GithubUsername, name) // 请求远端仓库最新 Tag 的 API
			goGiteaLatestSourceTagApi := color.Sprintf(general.GoLatestSourceTagApiFormat, config.Install.Go.GiteaApi, config.Install.Go.GiteaUsername, name)    // 请求远端仓库最新 Tag 的 API
			// 请求 API - GitHub
			body, err := general.RequestApi(goGithubLatestSourceTagApi)
			if err != nil {
				color.Error.Println(err)
				// 请求 API - Gitea
				body, err = general.RequestApi(goGiteaLatestSourceTagApi)
				if err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
			}
			// 获取远端版本（用于 source 安装方法）
			remoteTag, err := general.GetLatestSourceTag(body)
			if err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				continue
			}
			// 获取本地版本
			localProgram := filepath.Join(config.Install.ProgramPath, name) // 本地程序路径
			nameArgs := []string{"version", "--only"}                       // 本地程序参数
			localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
			// 比较远端和本地版本
			if remoteTag == localVersion { // 版本一致，则输出无需更新信息
				text := color.Sprintf("%s %s %s %s\n", general.LatestFlag, general.FgGreenText(name), general.FgYellowText(localVersion), general.FgWhiteText(general.LatestVersionMessage))
				color.Printf(text)
				textLength = general.RealLength(text) // 分隔符长度
			} else { // 版本不一致，则安装或更新程序，并输出已安装/更新信息
				// 如果 Temp 中已有远端仓库则删除重新克隆
				goSourceTempDir := filepath.Join(config.Install.SourceTemp, name)
				if general.FileExist(goSourceTempDir) {
					if err := os.RemoveAll(goSourceTempDir); err != nil {
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
				}
				// 克隆远端仓库 - GitHub
				goGithubCloneBaseUrl := color.Sprintf("%s/%s", config.Install.Go.GithubUrl, config.Install.Go.GithubUsername) // 远端仓库基础克隆地址（除仓库名）
				color.Printf("%s %s %s %s ", general.DownloadFlag, general.LightText("Clone"), general.FgGreenText(name), "from GitHub")
				if err := general.CloneRepoViaHTTP(config.Install.SourceTemp, goGithubCloneBaseUrl, name); err != nil {
					color.Printf("%s\n", general.ErrorText("error -> ", err))
					// 克隆远端仓库 - Gitea
					goGiteaCloneBaseUrl := color.Sprintf("%s/%s", config.Install.Go.GiteaUrl, config.Install.Go.GiteaUsername) // 远端仓库基础克隆地址（除仓库名）
					color.Printf("%s %s %s %s ", general.DownloadFlag, general.LightText("Clone"), general.FgGreenText(name), "from Gitea")
					if err := general.CloneRepoViaHTTP(config.Install.SourceTemp, goGiteaCloneBaseUrl, name); err != nil {
						text := color.Sprintf("%s\n", general.ErrorText("error -> ", err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					} else {
						color.Println(general.SuccessText("success"))
					}
				} else {
					color.Println(general.SuccessText("success"))
				}
				// 进到下载的远端文件目录
				if err := general.GoToDir(goSourceTempDir); err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
				// 编译生成程序
				if general.FileExist("Makefile") { // Makefile 文件存在则使用 make 编译
					makeArgs := []string{}
					if err := general.RunCommand("make", makeArgs); err != nil {
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
				} else if general.FileExist("main.go") { // Makefile 文件不存在则使用 `go build` 命令编译
					buildArgs := []string{"build", "-trimpath", "-ldflags=-s -w", "-o", name}
					if err := general.RunCommand("go", buildArgs); err != nil {
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
				} else {
					text := color.Error.Sprintf(general.UnableToCompileMessage)
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
				// 检测编译生成的程序是否存在
				compileProgram := filepath.Join(config.Install.SourceTemp, name, config.Install.Go.GeneratePath, name) // 编译生成的程序
				if general.FileExist(compileProgram) {
					// 检测本地程序是否存在
					if commandErr != nil { // 不存在，安装
						if general.FileExist("Makefile") { // Makefile 文件存在则使用 `make install` 命令安装
							makeArgs := []string{"install"}
							if err := general.RunCommand("make", makeArgs); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						} else { // Makefile 文件不存在则使用自定义函数安装
							if err := general.InstallFile(compileProgram, localProgram, 0755); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							} else {
								// 为已安装的程序设置可执行权限
								if err := os.Chmod(localProgram, 0755); err != nil {
									text := color.Sprintf("%s\n", general.ErrorText(err))
									color.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								}
							}
						}
						// 本次安装结束分隔符
						text := color.Sprintf("%s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(remoteTag), general.FgMagentaText("installed"))
						color.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					} else { // 存在，更新
						if general.FileExist("Makefile") { // Makefile 文件存在则使用 `make install` 命令更新
							makeArgs := []string{"install"}
							if err := general.RunCommand("make", makeArgs); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						} else { // Makefile 文件不存在则使用自定义函数更新
							if err := os.Remove(localProgram); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
							if err := general.InstallFile(compileProgram, localProgram, 0755); err != nil {
								text := color.Sprintf("%s\n", general.ErrorText(err))
								color.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							} else {
								// 为已安装的程序设置可执行权限
								if err := os.Chmod(localProgram, 0755); err != nil {
									text := color.Sprintf("%s\n", general.ErrorText(err))
									color.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								}
							}
						}
						// 本次更新结束分隔符
						text := color.Sprintf("%s %s %s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(localVersion), general.FgWhiteText("-->"), general.NoteText(remoteTag), general.FgMagentaText("updated"))
						color.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					}
					// 生成/更新自动补全脚本
					for _, completionDir := range config.Install.Go.CompletionDir {
						if general.FileExist(completionDir) {
							generateArgs := []string{"-c", color.Sprintf("%s completion zsh > %s/_%s", localProgram, completionDir, name)}
							if err := general.RunCommand("bash", generateArgs); err != nil {
								text := color.Sprintf("%s %s\n", general.ErrorFlag, general.ErrorText(general.AcsInstallFailedMessage))
								color.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
								continue
							} else {
								text := color.Sprintf("%s %s\n", general.SuccessFlag, general.SecondaryText(general.AcsInstallSuccessMessage))
								color.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
								break
							}
						}
					}
				} else {
					text := color.Error.Sprintf("The source file %s does not exist", compileProgram)
					color.Printf(text)
					textLength = general.RealLength(text) // 分隔符长度
				}
			}
			// 分隔符和延时（延时使输出更加顺畅）
			general.PrintDelimiter(textLength) // 分隔符
			general.Delay(0.1)                 // 0.01s
		}
	default:
		text := color.Error.Sprintf("Unsupported installation method '%s': only 'release' and 'source' are supported", config.Install.Method)
		color.Printf(text)
	}
}

// InstallShellBasedProgram 安装/更新基于 Shell 的程序
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
func InstallShellBasedProgram(configTree *toml.Tree) {
	// 获取配置项
	config, err := general.LoadConfigToStruct(configTree)
	if err != nil {
		color.Error.Println(err)
		return
	}
	// 开始安装
	color.Info.Tips("Install \x1b[3m%s\x1b[0m programs", general.FgCyanText("shell-based"))
	color.Info.Tips("%s: %s", general.FgWhiteText("Installation path"), general.PrimaryText(config.Install.ProgramPath))
	// 设置代理
	general.SetVariable("http_proxy", config.Variable.HTTPProxy)
	general.SetVariable("https_proxy", config.Variable.HTTPSProxy)

	// 设置进度条参数
	general.ProgressParameters["view"] = "0"

	// 创建临时目录
	if err := general.CreateDir(config.Install.SourceTemp); err != nil {
		color.Error.Println(err)
		return
	}

	// 让用户选择需要安装/更新的程序
	selectedNames, err := general.MultipleSelectionFilter(config.Install.Shell.Names)
	if err != nil {
		color.Error.Println(err)
	}
	// 对所选的程序进行排序
	sort.Strings(selectedNames)
	// 遍历所选脚本名
	for _, name := range selectedNames {
		textLength := 0                                                                                                                                                                                             // 输出文本的长度
		shellGithubLatestHashApi := color.Sprintf(general.ShellLatestHashApiFormat, config.Install.Shell.GithubApi, config.Install.Shell.GithubUsername, config.Install.Shell.Repo, config.Install.Shell.Dir, name) // 请求远端仓库最新脚本的 Hash 值的 API
		shellGiteaLatestHashApi := color.Sprintf(general.ShellLatestHashApiFormat, config.Install.Shell.GiteaApi, config.Install.Shell.GiteaUsername, config.Install.Shell.Repo, config.Install.Shell.Dir, name)    // 请求远端仓库最新脚本的 Hash 值的 API
		// 请求 API - GitHub
		body, err := general.RequestApi(shellGithubLatestHashApi)
		if err != nil {
			color.Error.Println(err)
			// 请求 API - Gitea
			body, err = general.RequestApi(shellGiteaLatestHashApi)
			if err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				continue
			}
		}
		// 获取远端脚本 Hash
		remoteHash, err := general.GetLatestSourceHash(body)
		if err != nil {
			text := color.Sprintf("%s\n", general.ErrorText(err))
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			continue
		}
		// 获取本地脚本 Hash
		localProgram := filepath.Join(config.Install.ProgramPath, name) // 本地程序路径
		gitHashObjectArgs := []string{"hash-object", localProgram}      // 本地程序参数
		localHash, commandErr := general.RunCommandGetResult("git", gitHashObjectArgs)
		// 比较远端和本地脚本 Hash
		if remoteHash == localHash { // Hash 一致，则输出无需更新信息
			text := color.Sprintf("%s %s %s\n", general.LatestFlag, general.FgGreenText(name), general.FgWhiteText(general.LatestVersionMessage))
			color.Printf(text)
			textLength = general.RealLength(text) // 分隔符长度
		} else { // Hash 不一致，则更新脚本，并输出已更新信息
			shellUrlFile := filepath.Join(config.Install.Shell.Dir, name)                                // 脚本在仓库中的路径
			scriptLocalPath := filepath.Join(config.Install.SourceTemp, config.Install.Shell.Repo, name) // 脚本本地存储位置
			// 下载远端脚本 - GitHub
			shellGithubBaseDownloadUrl := color.Sprintf(general.ShellGithubBaseDownloadUrlFormat, config.Install.Shell.GithubRaw, config.Install.Shell.GithubUsername, config.Install.Shell.Repo, config.Install.Shell.GithubBranch) // 脚本远端仓库基础地址
			fileUrl := color.Sprintf("%s/%s", shellGithubBaseDownloadUrl, shellUrlFile)
			if err := general.DownloadFile(fileUrl, scriptLocalPath, general.ProgressParameters); err != nil {
				color.Error.Println(err)
				// 下载远端脚本 - Gitea
				shellGiteaBaseDownloadUrl := color.Sprintf(general.ShellGiteaBaseDownloadUrlFormat, config.Install.Shell.GiteaRaw, config.Install.Shell.GiteaUsername, config.Install.Shell.Repo, config.Install.Shell.GiteaBranch) // 脚本远端仓库基础地址
				fileUrl := color.Sprintf("%s/%s", shellGiteaBaseDownloadUrl, shellUrlFile)
				if err = general.DownloadFile(fileUrl, scriptLocalPath, general.ProgressParameters); err != nil {
					text := color.Sprintf("%s\n", general.ErrorText(err))
					color.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
			}
			// 检测脚本文件是否存在
			if general.FileExist(scriptLocalPath) {
				// 检测本地程序是否存在
				if commandErr != nil { // 不存在，安装
					if err := general.InstallFile(scriptLocalPath, localProgram, 0755); err != nil {
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					} else {
						// 为已安装的脚本设置可执行权限
						if err := os.Chmod(localProgram, 0755); err != nil {
							color.Error.Println(err)
						}
						text := color.Sprintf("%s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(remoteHash[:6]), general.FgMagentaText("installed"))
						color.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					}
				} else { // 存在，更新
					if err := os.Remove(localProgram); err != nil {
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
					if err := general.InstallFile(scriptLocalPath, localProgram, 0755); err != nil {
						text := color.Sprintf("%s\n", general.ErrorText(err))
						color.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					} else {
						// 为已更新的脚本设置可执行权限
						if err := os.Chmod(localProgram, 0755); err != nil {
							color.Error.Println(err)
						}
						text := color.Sprintf("%s %s %s %s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgYellowText(localHash[:6]), general.FgWhiteText("-->"), general.NoteText(remoteHash[:6]), general.FgMagentaText("updated"))
						color.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					}
				}
			} else {
				text := color.Error.Sprintf("The source file %s does not exist\n", scriptLocalPath)
				color.Printf(text)
				textLength = general.RealLength(text) // 分隔符长度
			}
		}
		// 分隔符和延时（延时使输出更加顺畅）
		general.PrintDelimiter(textLength) // 分隔符
		general.Delay(0.1)                 // 0.1s
	}
}
