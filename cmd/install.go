/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 13:35:06

Description: 执行子命令 'install'
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yhyj/manager/cli"
	"github.com/yhyj/manager/general"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or update programs and scripts",
	Long:  `Install or update self-developed programs and scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		cfgFile, _ := cmd.Flags().GetString("config")
		allFlag, _ := cmd.Flags().GetBool("all")
		goFlag, _ := cmd.Flags().GetBool("go")
		shellFlag, _ := cmd.Flags().GetBool("shell")

		var (
			httpProxy  string // http 代理
			httpsProxy string // https 代理

			installMethod        string // 程序安装方法，目前支持 relrase 和 source
			installProgramPath   string // 程序安装路径
			installResourcesPath string // 资源安装路径
			installReleaseTemp   string // 使用 release 安装方法时下载文件的存储地址
			installSourceTemp    string // 使用 source 安装方法时下载文件的存储地址

			goNames          []interface{} // 要安装的基于 go 的程序名列表
			goReleaseApi     string        // 基于 go 的程序的 release 安装方法的 API 地址
			goGeneratePath   string        // 基于 go 的程序的 source 安装方法的编译生成路径
			goGithubUrl      string        // 基于 go 的程序的 source 安装方法的 GitHub 地址
			goGithubApi      string        // 基于 go 的程序的 source 安装方法的 GitHub API 地址
			goGithubUsername string        // 基于 go 的程序的 source 安装方法的 GitHub 用户名
			goGiteaUrl       string        // 基于 go 的程序的 source 安装方法的 Gitea 地址
			goGiteaApi       string        // 基于 go 的程序的 source 安装方法的 Gitea API 地址
			goGiteaUsername  string        // 基于 go 的程序的 source 安装方法的 Gitea 用户名
			goCompletionDir  []interface{} // 基于 go 的程序的自动补全文件安装目录列表

			shellNames          []interface{} // 要安装的 shell 脚本名列表
			shellGithubApi      string        // shell 脚本的 source 安装方法的 GitHub 地址
			shellGithubRaw      string        // shell 脚本的 source 安装方法的 GitHub 下载地址
			shellGithubBranch   string        // shell 脚本的 source 安装方法的 GitHub 分支名
			shellGithubUsername string        // shell 脚本的 source 安装方法的 GitHub 用户名
			shellGiteaApi       string        // shell 脚本的 source 安装方法的 Gitea 地址
			shellGiteaRaw       string        // shell 脚本的 source 安装方法的 Gitea 下载地址
			shellGiteaBranch    string        // shell 脚本的 source 安装方法的 Gitea 分支名
			shellGiteaUsername  string        // shell 脚本的 source 安装方法的 Gitea 用户名
			shellRepo           string        // shell 脚本所在的仓库名
			shellDir            string        // shell 脚本在仓库中的路径
		)

		// 检查配置文件是否存在
		configTree, err := cli.GetTomlConfig(cfgFile)
		if err != nil {
			fmt.Printf(general.ErrorBaseFormat, err)
			return
		} else {
			// 获取配置项
			if configTree.Has("install.method") {
				installMethod = configTree.Get("install.method").(string)
			}
			if configTree.Has("install.program_path") {
				installProgramPath = configTree.Get("install.program_path").(string)
			}
			if configTree.Has("install.resources_path") {
				installResourcesPath = configTree.Get("install.resources_path").(string)
			}
			if configTree.Has("install.release_temp") {
				installReleaseTemp = configTree.Get("install.release_temp").(string)
			}
			if configTree.Has("install.source_temp") {
				installSourceTemp = configTree.Get("install.source_temp").(string)
			}
			if configTree.Has("install.go.release_api") {
				goReleaseApi = configTree.Get("install.go.release_api").(string)
			}
			if configTree.Has("install.go.generate_path") {
				goGeneratePath = configTree.Get("install.go.generate_path").(string)
			}
			if configTree.Has("install.go.github_url") {
				goGithubUrl = configTree.Get("install.go.github_url").(string)
			}
			if configTree.Has("install.go.github_api") {
				goGithubApi = configTree.Get("install.go.github_api").(string)
			}
			if configTree.Has("install.go.github_username") {
				goGithubUsername = configTree.Get("install.go.github_username").(string)
			}
			if configTree.Has("install.go.gitea_url") {
				goGiteaUrl = configTree.Get("install.go.gitea_url").(string)
			}
			if configTree.Has("install.go.gitea_api") {
				goGiteaApi = configTree.Get("install.go.gitea_api").(string)
			}
			if configTree.Has("install.go.gitea_username") {
				goGiteaUsername = configTree.Get("install.go.gitea_username").(string)
			}
			if configTree.Has("install.go.names") {
				goNames = configTree.Get("install.go.names").([]interface{})
			}
			if configTree.Has("install.go.completion_dir") {
				goCompletionDir = configTree.Get("install.go.completion_dir").([]interface{})
			}
			if configTree.Has("install.shell.github_api") {
				shellGithubApi = configTree.Get("install.shell.github_api").(string)
			}
			if configTree.Has("install.shell.github_raw") {
				shellGithubRaw = configTree.Get("install.shell.github_raw").(string)
			}
			if configTree.Has("install.shell.github_branch") {
				shellGithubBranch = configTree.Get("install.shell.github_branch").(string)
			}
			if configTree.Has("install.shell.github_username") {
				shellGithubUsername = configTree.Get("install.shell.github_username").(string)
			}
			if configTree.Has("install.shell.gitea_api") {
				shellGiteaApi = configTree.Get("install.shell.gitea_api").(string)
			}
			if configTree.Has("install.shell.gitea_raw") {
				shellGiteaRaw = configTree.Get("install.shell.gitea_raw").(string)
			}
			if configTree.Has("install.shell.gitea_branch") {
				shellGiteaBranch = configTree.Get("install.shell.gitea_branch").(string)
			}
			if configTree.Has("install.shell.gitea_username") {
				shellGiteaUsername = configTree.Get("install.shell.gitea_username").(string)
			}
			if configTree.Has("install.shell.repo") {
				shellRepo = configTree.Get("install.shell.repo").(string)
			}
			if configTree.Has("install.shell.dir") {
				shellDir = configTree.Get("install.shell.dir").(string)
			}
			if configTree.Has("install.shell.names") {
				shellNames = configTree.Get("install.shell.names").([]interface{})
			}
			if configTree.Has("variable.http_proxy") {
				httpProxy = configTree.Get("variable.http_proxy").(string)
			}
			if configTree.Has("variable.https_proxy") {
				httpsProxy = configTree.Get("variable.https_proxy").(string)
			}
		}

		// 输出文本
		var (
			latestVersionMessage     = "is already the latest version"                 // 已安装的程序和脚本为最新版
			unableToCompileMessage   = "Makefile or main.go file does not exist"       // 缺失编译文件无法完成编译
			acsInstallSuccessMessage = "auto-completion script installed successfully" // 自动补全脚本安装成功
			acsInstallFailedMessage  = "auto-completion script installation failed"    // 自动补全脚本安装失败
		)

		// 字符串格式
		var (
			goLatestReleaseTagApiFormat      = "%s/repos/%s/%s/releases/latest" // 请求远端仓库最新 Tag 的 API - Release
			goLatestSourceTagApiFormat       = "%s/repos/%s/%s/tags"            // 请求远端仓库最新 Tag 的 API - Source
			shellLatestHashApiFormat         = "%s/repos/%s/%s/contents/%s/%s"  // 请求远端仓库最新脚本的 Hash 值的 API
			shellGithubBaseDownloadUrlFormat = "%s/%s/%s/%s"                    // 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - GitHub 格式
			shellGiteaBaseDownloadUrlFormat  = "%s/%s/%s/raw/branch/%s"         // 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - Gitea 格式
		)

		// 根据参数执行操作
		if allFlag {
			goFlag, shellFlag = true, true
		}

		// 安装/更新 shell 脚本
		if shellFlag {
			fmt.Printf(general.TitleH1Format, "Installing shell-based programs...")
			// 设置代理
			general.SetVariable("http_proxy", httpProxy)
			general.SetVariable("https_proxy", httpsProxy)
			// 创建临时目录
			if err := general.CreateDir(installSourceTemp); err != nil {
				fmt.Printf(general.ErrorBaseFormat, err)
				return
			}
			// 遍历所有脚本名
			for _, name := range shellNames {
				textLength := 0                                                                                                                            // 输出文本的长度
				shellGithubLatestHashApi := fmt.Sprintf(shellLatestHashApiFormat, shellGithubApi, shellGithubUsername, shellRepo, shellDir, name.(string)) // 请求远端仓库最新脚本的 Hash 值的 API
				shellGiteaLatestHashApi := fmt.Sprintf(shellLatestHashApiFormat, shellGiteaApi, shellGiteaUsername, shellRepo, shellDir, name.(string))    // 请求远端仓库最新脚本的 Hash 值的 API
				// 请求 API - GitHub
				body, err := general.RequestApi(shellGithubLatestHashApi)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					// 请求 API - Gitea
					body, err = general.RequestApi(shellGiteaLatestHashApi)
					if err != nil {
						text := fmt.Sprintf(general.ErrorBaseFormat, err)
						fmt.Printf(text)
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
					text := fmt.Sprintf(general.ErrorBaseFormat, err)
					fmt.Printf(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(0.1)                    // 0.1s
					continue
				}
				// 获取本地脚本 Hash
				localProgram := filepath.Join(installProgramPath, name.(string)) // 本地程序路径
				gitHashObjectArgs := []string{"hash-object", localProgram}       // 本地程序参数
				localHash, commandErr := general.RunCommandGetResult("git", gitHashObjectArgs)
				// 比较远端和本地脚本 Hash
				if remoteHash == localHash { // Hash 一致，则输出无需更新信息
					text := fmt.Sprintf(general.SliceTraverse2PSuffixFormat, general.Dot, " ", name.(string), " ", latestVersionMessage)
					fmt.Printf(text)
					textLength = general.RealLength(text) // 分隔符长度
				} else { // Hash 不一致，则更新脚本，并输出已更新信息
					shellUrlFile := filepath.Join(shellDir, name.(string))                        // 脚本在仓库中的路径
					scriptLocalPath := filepath.Join(installSourceTemp, shellRepo, name.(string)) // 脚本本地存储位置
					// 下载远端脚本 - GitHub
					shellGithubBaseDownloadUrl := fmt.Sprintf(shellGithubBaseDownloadUrlFormat, shellGithubRaw, shellGithubUsername, shellRepo, shellGithubBranch) // 脚本远端仓库基础地址
					fileUrl := fmt.Sprintf("%s/%s", shellGithubBaseDownloadUrl, shellUrlFile)
					if err := cli.DownloadFile(fileUrl, scriptLocalPath, general.ProgressParameters); err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						// 下载远端脚本 - Gitea
						shellGiteaBaseDownloadUrl := fmt.Sprintf(shellGiteaBaseDownloadUrlFormat, shellGiteaRaw, shellGiteaUsername, shellRepo, shellGiteaBranch) // 脚本远端仓库基础地址
						fileUrl := fmt.Sprintf("%s/%s", shellGiteaBaseDownloadUrl, shellUrlFile)
						if err = cli.DownloadFile(fileUrl, scriptLocalPath, general.ProgressParameters); err != nil {
							text := fmt.Sprintf(general.ErrorBaseFormat, err)
							fmt.Printf(text)
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
							if err := cli.InstallFile(scriptLocalPath, localProgram, 0755); err != nil {
								text := fmt.Sprintf(general.ErrorBaseFormat, err)
								fmt.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							} else {
								// 为已安装的脚本设置可执行权限
								if err := os.Chmod(localProgram, 0755); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
								}
								text := fmt.Sprintf(general.SliceTraverse4PFormat, general.Yes, " ", name.(string), " ", remoteHash[:6], " ", "installed")
								fmt.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
							}
						} else { // 存在，更新
							if err := os.Remove(localProgram); err != nil {
								text := fmt.Sprintf(general.ErrorBaseFormat, err)
								fmt.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
							if err := cli.InstallFile(scriptLocalPath, localProgram, 0755); err != nil {
								text := fmt.Sprintf(general.ErrorBaseFormat, err)
								fmt.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							} else {
								// 为已更新的脚本设置可执行权限
								if err := os.Chmod(localProgram, 0755); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
								}
								text := fmt.Sprintf(general.SliceTraverse5PFormat, general.Yes, " ", name.(string), " ", localHash[:6], " -> ", remoteHash[:6], " ", "updated")
								fmt.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
							}
						}
					} else {
						text := fmt.Sprintf(general.ErrorBaseFormat, fmt.Sprintf("The source file %s does not exist", scriptLocalPath))
						fmt.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					}
				}
				// 分隔符和延时（延时使输出更加顺畅）
				general.PrintDelimiter(textLength) // 分隔符
				general.Delay(0.1)                 // 0.1s
			}
		}

		// 安装/更新基于 go 的程序
		if goFlag {
			fmt.Printf(general.TitleH1Format, "Installing go-based programs...")
			// 设置代理
			general.SetVariable("http_proxy", httpProxy)
			general.SetVariable("https_proxy", httpsProxy)

			// 设置进度条参数
			general.ProgressParameters["view"] = "1"
			general.ProgressParameters["sep"] = "-"

			// 使用配置的安装方式进行安装
			if strings.ToLower(installMethod) == "release" {
				// 创建临时目录
				if err := general.CreateDir(installReleaseTemp); err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					return
				}
				// 遍历所有程序名
				for _, name := range goNames {
					textLength := 0                                                                                                        // 输出文本的长度
					goGithubLatestReleaseTagApi := fmt.Sprintf(goLatestReleaseTagApiFormat, goReleaseApi, goGithubUsername, name.(string)) // 请求远端仓库最新 Tag 的 API
					// 请求 API - GitHub
					body, err := general.RequestApi(goGithubLatestReleaseTagApi)
					if err != nil {
						text := fmt.Sprintf(general.ErrorBaseFormat, err)
						fmt.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
					// 获取远端版本（用于 release 安装方法）
					remoteTag, err := general.GetLatestReleaseTag(body)
					if err != nil {
						text := fmt.Sprintf(general.ErrorBaseFormat, err)
						fmt.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
					// 获取本地版本
					localProgram := filepath.Join(installProgramPath, name.(string)) // 本地程序路径
					nameArgs := []string{"version", "--only"}                        // 本地程序参数
					localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
					// 比较远端和本地版本
					if remoteTag == localVersion { // 版本一致，则输出无需更新信息
						text := fmt.Sprintf(general.SliceTraverse3PSuffixFormat, general.Dot, " ", name.(string), " ", remoteTag, " ", latestVersionMessage)
						fmt.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					} else { // 版本不一致，则安装或更新程序，并输出已安装/更新信息
						// 下载远端文件（如果 Temp 中已有远端文件则删除重新下载）
						goReleaseTempDir := filepath.Join(installReleaseTemp, name.(string))
						if general.FileExist(goReleaseTempDir) {
							if err := os.RemoveAll(goReleaseTempDir); err != nil {
								text := fmt.Sprintf(general.ErrorBaseFormat, err)
								fmt.Printf(text)
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
						archiveFileNameWithoutFileType := fmt.Sprintf("%s_%s_%s_%s", name.(string), remoteTag, general.Platform, general.Arch)
						fileName.ArchiveFile = fmt.Sprintf("%s.%s", archiveFileNameWithoutFileType, fileType)
						// 获取 Release 文件信息
						filesInfo, err := general.GetReleaseFileInfo(body, fileName)
						if err != nil {
							text := fmt.Sprintf(general.ErrorBaseFormat, err)
							fmt.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							continue
						}
						// fmt.Printf(general.SliceTraverse2PSuffixFormat, general.Run, " Download ", fmt.Sprintf("[%s] - %s", name, filesInfo.ChecksumsFileInfo.Name), " ", "from GitHub Release ")
						general.ProgressParameters["action"] = general.Run
						general.ProgressParameters["prefix"] = "Download"
						general.ProgressParameters["project"] = fmt.Sprintf("[%s]", name)
						general.ProgressParameters["fileName"] = fmt.Sprintf("[%s]", filesInfo.ChecksumsFileInfo.Name)
						general.ProgressParameters["suffix"] = "from Github release:"
						checksumsLocalPath := filepath.Join(installReleaseTemp, name.(string), filesInfo.ChecksumsFileInfo.Name) // Checksums 文件本地存储位置
						if err := cli.DownloadFile(filesInfo.ChecksumsFileInfo.DownloadUrl, checksumsLocalPath, general.ProgressParameters); err != nil {
							text := fmt.Sprintf(general.ErrorSuffixFormat, "error", " -> ", err)
							fmt.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							continue
						}
						// fmt.Printf(general.SliceTraverse2PSuffixFormat, general.Run, " Download ", fmt.Sprintf("[%s] - %s", name, filesInfo.ArchiveFileInfo.Name), " ", "from GitHub Release ")
						general.ProgressParameters["action"] = general.Run
						general.ProgressParameters["prefix"] = "Download"
						general.ProgressParameters["project"] = fmt.Sprintf("[%s]", name)
						general.ProgressParameters["fileName"] = fmt.Sprintf("[%s]", filesInfo.ArchiveFileInfo.Name)
						general.ProgressParameters["suffix"] = "from Github release:"
						archiveLocalPath := filepath.Join(installReleaseTemp, name.(string), filesInfo.ArchiveFileInfo.Name) // Release 文件本地存储位置
						if err := cli.DownloadFile(filesInfo.ArchiveFileInfo.DownloadUrl, archiveLocalPath, general.ProgressParameters); err != nil {
							text := fmt.Sprintf(general.ErrorSuffixFormat, "error", " -> ", err)
							fmt.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							continue
						}
						// 进到下载的远端文件目录
						if err := general.GoToDir(goReleaseTempDir); err != nil {
							text := fmt.Sprintf(general.ErrorBaseFormat, err)
							fmt.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							continue
						}
						// 使用校验文件校验下载的压缩包
						verificationResult, err := cli.FileVerification(checksumsLocalPath, archiveLocalPath)
						if err != nil {
							text := fmt.Sprintf(general.ErrorBaseFormat, err)
							fmt.Printf(text)
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
								text := fmt.Sprintf(general.ErrorBaseFormat, err)
								fmt.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
							archivedProgram := filepath.Join(goReleaseTempDir, archiveFileNameWithoutFileType, name.(string))       // 解压得到的程序
							archivedResourcesFolder := filepath.Join(goReleaseTempDir, archiveFileNameWithoutFileType, "resources") // 解压得到的资源文件夹
							// 检测本地程序是否存在
							if commandErr != nil { // 不存在，安装
								if err := cli.InstallFile(archivedProgram, localProgram, 0755); err != nil { // 安装程序
									text := fmt.Sprintf(general.ErrorBaseFormat, err)
									fmt.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								} else { // 为已安装的程序设置可执行权限
									if err := os.Chmod(localProgram, 0755); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
								}
								// 安装资源文件 - desktop 文件
								archivedResourcesDesktopFile := filepath.Join(archivedResourcesFolder, "applications", fmt.Sprintf("%s.desktop", name.(string))) // 解压得到的资源文件 - desktop 文件
								localResourcesDesktopFile := filepath.Join(installResourcesPath, "applications", fmt.Sprintf("%s.desktop", name.(string)))       // 本地资源文件 - desktop 文件
								if general.FileExist(archivedResourcesDesktopFile) {
									if err := cli.InstallFile(archivedResourcesDesktopFile, localResourcesDesktopFile, 0644); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
								}
								// 安装资源文件 - icon 文件
								archivedResourcesIconFolder := filepath.Join(archivedResourcesFolder, "pixmaps") // 解压得到的资源文件 - icon 文件夹
								localResourcesIconFolder := filepath.Join(installResourcesPath, "pixmaps")       // 本地资源文件 - icon 文件夹
								if general.FileExist(archivedResourcesIconFolder) {
									files, err := general.ListFolderFiles(archivedResourcesIconFolder)
									if err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
									if !general.FileExist(localResourcesIconFolder) {
										err := general.CreateDir(localResourcesIconFolder)
										if err != nil {
											text := fmt.Sprintf(general.ErrorBaseFormat, err)
											fmt.Printf(text)
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
										if err := cli.InstallFile(archivedResourcesIconFile, localResourcesIconFile, 0644); err != nil {
											text := fmt.Sprintf(general.ErrorBaseFormat, err)
											fmt.Printf(text)
											// 分隔符和延时（延时使输出更加顺畅）
											textLength = general.RealLength(text) // 分隔符长度
											general.PrintDelimiter(textLength)    // 分隔符
											general.Delay(0.1)                    // 0.1s
											continue
										}
									}
								}
								// 本次安装结束分隔符
								text := fmt.Sprintf(general.SliceTraverse4PFormat, general.Yes, " ", name.(string), " ", remoteTag, " ", "installed")
								fmt.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
							} else { // 存在，更新
								if err := os.Remove(localProgram); err != nil { // 删除已安装的旧程序
									text := fmt.Sprintf(general.ErrorBaseFormat, err)
									fmt.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								}
								if err := cli.InstallFile(archivedProgram, localProgram, 0755); err != nil { // 安装新程序
									text := fmt.Sprintf(general.ErrorBaseFormat, err)
									fmt.Printf(text)
									// 分隔符和延时（延时使输出更加顺畅）
									textLength = general.RealLength(text) // 分隔符长度
									general.PrintDelimiter(textLength)    // 分隔符
									general.Delay(0.1)                    // 0.1s
									continue
								} else { // 为已安装的程序设置可执行权限
									if err := os.Chmod(localProgram, 0755); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
								}
								// 安装资源文件 - desktop 文件
								archivedResourcesDesktopFile := filepath.Join(archivedResourcesFolder, "applications", fmt.Sprintf("%s.desktop", name.(string))) // 解压得到的资源文件 - desktop 文件
								localResourcesDesktopFile := filepath.Join(installResourcesPath, "applications", fmt.Sprintf("%s.desktop", name.(string)))       // 本地资源文件 - desktop 文件
								if general.FileExist(archivedResourcesDesktopFile) {
									if err := cli.InstallFile(archivedResourcesDesktopFile, localResourcesDesktopFile, 0644); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
								}
								// 安装资源文件 - icon 文件
								archivedResourcesIconFolder := filepath.Join(archivedResourcesFolder, "pixmaps") // 解压得到的资源文件 - icon 文件夹
								localResourcesIconFolder := filepath.Join(installResourcesPath, "pixmaps")       // 本地资源文件 - icon 文件夹
								if general.FileExist(archivedResourcesIconFolder) {
									files, err := general.ListFolderFiles(archivedResourcesIconFolder)
									if err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
									if !general.FileExist(localResourcesIconFolder) {
										err := general.CreateDir(localResourcesIconFolder)
										if err != nil {
											text := fmt.Sprintf(general.ErrorBaseFormat, err)
											fmt.Printf(text)
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
										if err := cli.InstallFile(archivedResourcesIconFile, localResourcesIconFile, 0644); err != nil {
											text := fmt.Sprintf(general.ErrorBaseFormat, err)
											fmt.Printf(text)
											// 分隔符和延时（延时使输出更加顺畅）
											textLength = general.RealLength(text) // 分隔符长度
											general.PrintDelimiter(textLength)    // 分隔符
											general.Delay(0.1)                    // 0.1s
											continue
										}
									}
								}
								// 本次更新结束分隔符
								text := fmt.Sprintf(general.SliceTraverse5PFormat, general.Yes, " ", name.(string), " ", localVersion, " -> ", remoteTag, " ", "updated")
								fmt.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
							}
							// 生成/更新自动补全脚本
							for _, completionDir := range goCompletionDir {
								if general.FileExist(completionDir.(string)) {
									generateArgs := []string{"-c", fmt.Sprintf("%s completion zsh > %s/_%s", localProgram, completionDir.(string), name.(string))}
									if err := general.RunCommand("bash", generateArgs); err != nil {
										text := fmt.Sprintf(general.ErrorSuffixFormat, general.No, " ", acsInstallFailedMessage)
										fmt.Printf(text)
										textLength = general.RealLength(text) // 分隔符长度
										continue
									} else {
										text := fmt.Sprintf(general.SuccessSuffixFormat, general.Yes, " ", acsInstallSuccessMessage)
										fmt.Printf(text)
										textLength = general.RealLength(text) // 分隔符长度
										break
									}
								}
							}
						} else { // 压缩包校验失败
							text := fmt.Sprintf(general.ErrorSuffixFormat, "Archive file verification failed", ": ", filesInfo.ArchiveFileInfo.Name)
							fmt.Printf(text)
							textLength = general.RealLength(text) // 分隔符长度
						}
					}
					// 分隔符和延时（延时使输出更加顺畅）
					general.PrintDelimiter(textLength) // 分隔符
					general.Delay(0.1)                 // 0.01s
				}
			} else if strings.ToLower(installMethod) == "source" {
				// 创建临时目录
				if err := general.CreateDir(installSourceTemp); err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					return
				}
				// 遍历所有程序名
				for _, name := range goNames {
					textLength := 0                                                                                                     // 输出文本的长度
					goGithubLatestSourceTagApi := fmt.Sprintf(goLatestSourceTagApiFormat, goGithubApi, goGithubUsername, name.(string)) // 请求远端仓库最新 Tag 的 API
					goGiteaLatestSourceTagApi := fmt.Sprintf(goLatestSourceTagApiFormat, goGiteaApi, goGiteaUsername, name.(string))    // 请求远端仓库最新 Tag 的 API
					// 请求 API - GitHub
					body, err := general.RequestApi(goGithubLatestSourceTagApi)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						// 请求 API - Gitea
						body, err = general.RequestApi(goGiteaLatestSourceTagApi)
						if err != nil {
							text := fmt.Sprintf(general.ErrorBaseFormat, err)
							fmt.Printf(text)
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
						text := fmt.Sprintf(general.ErrorBaseFormat, err)
						fmt.Printf(text)
						// 分隔符和延时（延时使输出更加顺畅）
						textLength = general.RealLength(text) // 分隔符长度
						general.PrintDelimiter(textLength)    // 分隔符
						general.Delay(0.1)                    // 0.1s
						continue
					}
					// 获取本地版本
					localProgram := filepath.Join(installProgramPath, name.(string)) // 本地程序路径
					nameArgs := []string{"version", "--only"}                        // 本地程序参数
					localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
					// 比较远端和本地版本
					if remoteTag == localVersion { // 版本一致，则输出无需更新信息
						text := fmt.Sprintf(general.SliceTraverse3PSuffixFormat, general.Dot, " ", name.(string), " ", remoteTag, " ", latestVersionMessage)
						fmt.Printf(text)
						textLength = general.RealLength(text) // 分隔符长度
					} else { // 版本不一致，则安装或更新程序，并输出已安装/更新信息
						// 如果 Temp 中已有远端仓库则删除重新克隆
						goSourceTempDir := filepath.Join(installSourceTemp, name.(string))
						if general.FileExist(goSourceTempDir) {
							if err := os.RemoveAll(goSourceTempDir); err != nil {
								text := fmt.Sprintf(general.ErrorBaseFormat, err)
								fmt.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						}
						// 克隆远端仓库 - GitHub
						goGithubCloneBaseUrl := fmt.Sprintf("%s/%s", goGithubUrl, goGithubUsername) // 远端仓库基础克隆地址（除仓库名）
						fmt.Printf(general.SliceTraverse2PSuffixFormat, general.Run, " Clone ", name.(string), " ", "from GitHub ")
						if err := cli.CloneRepoViaHTTP(installSourceTemp, goGithubCloneBaseUrl, name.(string)); err != nil {
							fmt.Printf(general.ErrorSuffixFormat, "error", " -> ", err)
							// 克隆远端仓库 - Gitea
							goGiteaCloneBaseUrl := fmt.Sprintf("%s/%s", goGiteaUrl, goGiteaUsername) // 远端仓库基础克隆地址（除仓库名）
							fmt.Printf(general.SliceTraverse2PSuffixFormat, general.Run, " Clone ", name.(string), " ", "from Gitea ")
							if err := cli.CloneRepoViaHTTP(installSourceTemp, goGiteaCloneBaseUrl, name.(string)); err != nil {
								text := fmt.Sprintf(general.ErrorSuffixFormat, "error", " -> ", err)
								fmt.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							} else {
								fmt.Printf(general.SuccessFormat, "success")
							}
						} else {
							fmt.Printf(general.SuccessFormat, "success")
						}
						// 进到下载的远端文件目录
						if err := general.GoToDir(goSourceTempDir); err != nil {
							text := fmt.Sprintf(general.ErrorBaseFormat, err)
							fmt.Printf(text)
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
								text := fmt.Sprintf(general.ErrorBaseFormat, err)
								fmt.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						} else if general.FileExist("main.go") { // Makefile 文件不存在则使用 `go build` 命令编译
							buildArgs := []string{"build", "-trimpath", "-ldflags=-s -w", "-o", name.(string)}
							if err := general.RunCommand("go", buildArgs); err != nil {
								text := fmt.Sprintf(general.ErrorBaseFormat, err)
								fmt.Printf(text)
								// 分隔符和延时（延时使输出更加顺畅）
								textLength = general.RealLength(text) // 分隔符长度
								general.PrintDelimiter(textLength)    // 分隔符
								general.Delay(0.1)                    // 0.1s
								continue
							}
						} else {
							text := fmt.Sprintf(general.ErrorBaseFormat, unableToCompileMessage)
							fmt.Printf(text)
							// 分隔符和延时（延时使输出更加顺畅）
							textLength = general.RealLength(text) // 分隔符长度
							general.PrintDelimiter(textLength)    // 分隔符
							general.Delay(0.1)                    // 0.1s
							continue
						}
						// 检测编译生成的程序是否存在
						compileProgram := filepath.Join(installSourceTemp, name.(string), goGeneratePath, name.(string)) // 编译生成的程序
						if general.FileExist(compileProgram) {
							// 检测本地程序是否存在
							if commandErr != nil { // 不存在，安装
								if general.FileExist("Makefile") { // Makefile 文件存在则使用 `make install` 命令安装
									makeArgs := []string{"install"}
									if err := general.RunCommand("make", makeArgs); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
								} else { // Makefile 文件不存在则使用自定义函数安装
									if err := cli.InstallFile(compileProgram, localProgram, 0755); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									} else {
										// 为已安装的程序设置可执行权限
										if err := os.Chmod(localProgram, 0755); err != nil {
											text := fmt.Sprintf(general.ErrorBaseFormat, err)
											fmt.Printf(text)
											// 分隔符和延时（延时使输出更加顺畅）
											textLength = general.RealLength(text) // 分隔符长度
											general.PrintDelimiter(textLength)    // 分隔符
											general.Delay(0.1)                    // 0.1s
											continue
										}
									}
								}
								// 本次安装结束分隔符
								text := fmt.Sprintf(general.SliceTraverse4PFormat, general.Yes, " ", name.(string), " ", remoteTag, " ", "installed")
								fmt.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
							} else { // 存在，更新
								if general.FileExist("Makefile") { // Makefile 文件存在则使用 `make install` 命令更新
									makeArgs := []string{"install"}
									if err := general.RunCommand("make", makeArgs); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
								} else { // Makefile 文件不存在则使用自定义函数更新
									if err := os.Remove(localProgram); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									}
									if err := cli.InstallFile(compileProgram, localProgram, 0755); err != nil {
										text := fmt.Sprintf(general.ErrorBaseFormat, err)
										fmt.Printf(text)
										// 分隔符和延时（延时使输出更加顺畅）
										textLength = general.RealLength(text) // 分隔符长度
										general.PrintDelimiter(textLength)    // 分隔符
										general.Delay(0.1)                    // 0.1s
										continue
									} else {
										// 为已安装的程序设置可执行权限
										if err := os.Chmod(localProgram, 0755); err != nil {
											text := fmt.Sprintf(general.ErrorBaseFormat, err)
											fmt.Printf(text)
											// 分隔符和延时（延时使输出更加顺畅）
											textLength = general.RealLength(text) // 分隔符长度
											general.PrintDelimiter(textLength)    // 分隔符
											general.Delay(0.1)                    // 0.1s
											continue
										}
									}
								}
								// 本次更新结束分隔符
								text := fmt.Sprintf(general.SliceTraverse5PFormat, general.Yes, " ", name.(string), " ", localVersion, " -> ", remoteTag, " ", "updated")
								fmt.Printf(text)
								textLength = general.RealLength(text) // 分隔符长度
							}
							// 生成/更新自动补全脚本
							for _, completionDir := range goCompletionDir {
								if general.FileExist(completionDir.(string)) {
									generateArgs := []string{"-c", fmt.Sprintf("%s completion zsh > %s/_%s", localProgram, completionDir.(string), name.(string))}
									if err := general.RunCommand("bash", generateArgs); err != nil {
										text := fmt.Sprintf(general.ErrorSuffixFormat, general.No, " ", acsInstallFailedMessage)
										fmt.Printf(text)
										textLength = general.RealLength(text) // 分隔符长度
										continue
									} else {
										text := fmt.Sprintf(general.SuccessSuffixFormat, general.Yes, " ", acsInstallSuccessMessage)
										fmt.Printf(text)
										textLength = general.RealLength(text) // 分隔符长度
										break
									}
								}
							}
						} else {
							text := fmt.Sprintf(general.ErrorBaseFormat, fmt.Sprintf("The source file %s does not exist", compileProgram))
							fmt.Printf(text)
							textLength = general.RealLength(text) // 分隔符长度
						}
					}
					// 分隔符和延时（延时使输出更加顺畅）
					general.PrintDelimiter(textLength) // 分隔符
					general.Delay(0.1)                 // 0.01s
				}
			} else {
				text := fmt.Sprintf(general.ErrorSuffixFormat, fmt.Sprintf("Unsupported installation method '%s'", installMethod), ": ", "only 'release' and 'source' are supported")
				fmt.Printf(text)
			}
		}
	},
}

func init() {
	installCmd.Flags().BoolP("all", "", false, "Install or update all programs and scripts")
	installCmd.Flags().BoolP("go", "", false, "Install or update programs developed based on go")
	installCmd.Flags().BoolP("shell", "", false, "Install or update shell scripts")

	installCmd.Flags().BoolP("help", "h", false, "help for install command")
	rootCmd.AddCommand(installCmd)
}
