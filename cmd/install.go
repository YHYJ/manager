/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 13:35:06

Description: 由程序子命令 install 执行
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

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

			installMethod      string // 程序安装方法，目前支持 relrase 和 source
			installPath        string // 程序安装路径
			installReleaseTemp string // 使用 release 安装方法时下载文件的存储地址
			installSourceTemp  string // 使用 source 安装方法时下载文件的存储地址

			goNames          []interface{} // 要安装的基于 go 的程序名列表
			goReleaseApi     string        // 基于 go 的程序的 release 安装方法的 API 地址
			goGeneratePath   string        // 基于 go 的程序的 source 安装方法的编译生成路径
			goGithubUrl      string        // 基于 go 的程序的 source 安装方法的 Github 地址
			goGithubApi      string        // 基于 go 的程序的 source 安装方法的 Github API 地址
			goGithubUsername string        // 基于 go 的程序的 source 安装方法的 Github 用户名
			goGiteaUrl       string        // 基于 go 的程序的 source 安装方法的 Gitea 地址
			goGiteaApi       string        // 基于 go 的程序的 source 安装方法的 Gitea API 地址
			goGiteaUsername  string        // 基于 go 的程序的 source 安装方法的 Gitea 用户名
			goCompletionDir  []interface{} // 基于 go 的程序的自动补全文件安装目录列表

			shellNames          []interface{} // 要安装的 shell 脚本名列表
			shellGithubApi      string        // shell 脚本的 source 安装方法的 Github 地址
			shellGithubRaw      string        // shell 脚本的 source 安装方法的 Github 下载地址
			shellGithubBranch   string        // shell 脚本的 source 安装方法的 Github 分支名
			shellGithubUsername string        // shell 脚本的 source 安装方法的 Github 用户名
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
			if configTree.Has("install.path") {
				installPath = configTree.Get("install.path").(string)
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
			shellGithubBaseDownloadUrlFormat = "%s/%s/%s/%s"                    // 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - Github 格式
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
				// 请求 API
				body, err := general.RequestApi(shellGithubLatestHashApi)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					body, err = general.RequestApi(shellGiteaLatestHashApi)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						continue
					}
				}
				// 获取远端脚本 Hash
				remoteHash, err := general.GetLatestSourceHash(body)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					continue
				}
				// 获取本地脚本 Hash
				localProgram := filepath.Join(installPath, name.(string))  // 本地程序路径
				gitHashObjectArgs := []string{"hash-object", localProgram} // 本地程序参数
				localHash, commandErr := general.RunCommandGetResult("git", gitHashObjectArgs)
				// 比较远端和本地脚本 Hash
				if remoteHash == localHash { // Hash 一致，则输出无需更新信息
					text := fmt.Sprintf(general.SliceTraverse2PSuffixFormat, "==>", " ", name.(string), " ", latestVersionMessage)
					fmt.Printf(text)
					controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
					textLength = len(controlRegex.ReplaceAllString(text, ""))
				} else { // Hash 不一致，则更新脚本，并输出已更新信息
					// 下载远端脚本
					shellGithubBaseDownloadUrl := fmt.Sprintf(shellGithubBaseDownloadUrlFormat, shellGithubRaw, shellGithubUsername, shellRepo, shellGithubBranch) // 脚本远端仓库基础地址
					shellGiteaBaseDownloadUrl := fmt.Sprintf(shellGiteaBaseDownloadUrlFormat, shellGiteaRaw, shellGiteaUsername, shellRepo, shellGiteaBranch)      // 脚本远端仓库基础地址
					shellUrlFile := filepath.Join(shellDir, name.(string))                                                                                         // 脚本在仓库中的路径
					scriptLocalPath := filepath.Join(installSourceTemp, shellRepo, name.(string))                                                                  // 脚本本地存储位置
					fileUrl := fmt.Sprintf("%s/%s", shellGithubBaseDownloadUrl, shellUrlFile)
					if err := cli.DownloadFile(fileUrl, scriptLocalPath); err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						fileUrl := fmt.Sprintf("%s/%s", shellGiteaBaseDownloadUrl, shellUrlFile)
						if err = cli.DownloadFile(fileUrl, scriptLocalPath); err != nil {
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
					}
					// 检测脚本文件是否存在
					if general.FileExist(scriptLocalPath) {
						// 检测本地程序是否存在
						if commandErr != nil { // 不存在，安装
							if err := cli.InstallFile(scriptLocalPath, localProgram, 0755); err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
								continue
							} else {
								// 为已安装的脚本设置可执行权限
								if err := os.Chmod(localProgram, 0755); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
								}
								text := fmt.Sprintf(general.SliceTraverse4PFormat, "==>", " ", name.(string), " ", remoteHash[:4], " ", "installed")
								fmt.Printf(text)
								controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
								textLength = len(controlRegex.ReplaceAllString(text, ""))
							}
						} else { // 存在，更新
							if err := os.Remove(localProgram); err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
							}
							if err := cli.InstallFile(scriptLocalPath, localProgram, 0755); err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
								continue
							} else {
								// 为已更新的脚本设置可执行权限
								if err := os.Chmod(localProgram, 0755); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
								}
								text := fmt.Sprintf(general.SliceTraverse5PFormat, "==>", " ", name.(string), " ", localHash, " -> ", remoteHash[:4], " ", "updated")
								fmt.Printf(text)
								controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
								textLength = len(controlRegex.ReplaceAllString(text, ""))
							}
						}
					} else {
						text := fmt.Sprintf(general.ErrorBaseFormat, fmt.Sprintf("The source file %s does not exist", scriptLocalPath))
						fmt.Printf(text)
						controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
						textLength = len(controlRegex.ReplaceAllString(text, ""))
					}
				}
				// 分隔符和延时
				dashes := strings.Repeat("-", textLength-1)  //组装分隔符（减去行尾换行符的一个长度）
				fmt.Printf(general.LineHiddenFormat, dashes) // 美化输出
				time.Sleep(100 * time.Millisecond)           // 添加一个0.01秒的延时，使输出更加顺畅
			}
		}

		// 安装/更新基于 go 的程序
		if goFlag {
			fmt.Printf(general.TitleH1Format, "Installing go-based programs...")
			// 设置代理
			general.SetVariable("http_proxy", httpProxy)
			general.SetVariable("https_proxy", httpsProxy)

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
					// 请求 API
					body, err := general.RequestApi(goGithubLatestReleaseTagApi)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						continue
					}
					// 获取远端版本（用于 release 安装方法）
					remoteTag, err := general.GetLatestReleaseTag(body)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						continue
					}
					// 获取本地版本
					localProgram := filepath.Join(installPath, name.(string)) // 本地程序路径
					nameArgs := []string{"version", "--only"}                 // 本地程序参数
					localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
					// 比较远端和本地版本
					if remoteTag == localVersion { // 版本一致，则输出无需更新信息
						text := fmt.Sprintf(general.SliceTraverse3PSuffixFormat, "==>", " ", name.(string), " ", remoteTag, " ", latestVersionMessage)
						fmt.Printf(text)
						controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
						textLength = len(controlRegex.ReplaceAllString(text, ""))
					} else { // 版本不一致，则安装或更新程序，并输出已安装/更新信息
						// 下载远端文件（如果 Temp 中已有远端文件则删除重新下载）
						goReleaseTempDir := filepath.Join(installReleaseTemp, name.(string))
						if general.FileExist(goReleaseTempDir) {
							if err := os.RemoveAll(goReleaseTempDir); err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
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
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
						fmt.Printf(general.SliceTraverse2PSuffixNoNewLineFormat, "==>", " Download ", filesInfo.ChecksumsFileInfo.Name, " ", "from GitHub Release ")
						checksumsLocalPath := filepath.Join(installReleaseTemp, name.(string), filesInfo.ChecksumsFileInfo.Name) // Checksums 文件本地存储位置
						if err := cli.DownloadFile(filesInfo.ChecksumsFileInfo.DownloadUrl, checksumsLocalPath); err != nil {
							fmt.Printf(general.ErrorSuffixFormat, "error", " -> ", err)
							continue
						} else {
							fmt.Printf(general.SuccessFormat, "success")
						}
						fmt.Printf(general.SliceTraverse2PSuffixNoNewLineFormat, "==>", " Download ", filesInfo.ArchiveFileInfo.Name, " ", "from GitHub Release ")
						archiveLocalPath := filepath.Join(installReleaseTemp, name.(string), filesInfo.ArchiveFileInfo.Name) // Release 文件本地存储位置
						if err := cli.DownloadFile(filesInfo.ArchiveFileInfo.DownloadUrl, archiveLocalPath); err != nil {
							fmt.Printf(general.ErrorSuffixFormat, "error", " -> ", err)
							continue
						} else {
							fmt.Printf(general.SuccessFormat, "success")
						}
						// 进到下载的远端文件目录
						if err := general.GoToDir(goReleaseTempDir); err != nil {
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
						// 使用校验文件校验下载的压缩包
						verificationResult, err := cli.FileVerification(checksumsLocalPath, archiveLocalPath)
						if err != nil {
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
						if verificationResult { // 压缩包校验通过
							// 解压压缩包
							err := general.UnzipFile(archiveLocalPath, goReleaseTempDir)
							if err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
								continue
							}
							archivedProgram := filepath.Join(goReleaseTempDir, archiveFileNameWithoutFileType, name.(string)) // 解压得到的程序
							// 检测本地程序是否存在
							if commandErr != nil { // 不存在，安装
								if err := cli.InstallFile(archivedProgram, localProgram, 0755); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
									continue
								} else {
									// 为已安装的脚本设置可执行权限
									if err := os.Chmod(localProgram, 0755); err != nil {
										fmt.Printf(general.ErrorBaseFormat, err)
									}
								}
								text := fmt.Sprintf(general.SliceTraverse4PFormat, "==>", " ", name.(string), " ", remoteTag, " ", "installed")
								fmt.Printf(text)
								controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
								textLength = len(controlRegex.ReplaceAllString(text, ""))
							} else { // 存在，更新
								if err := os.Remove(localProgram); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
								}
								if err := cli.InstallFile(archivedProgram, localProgram, 0755); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
									continue
								} else {
									// 为已安装的脚本设置可执行权限
									if err := os.Chmod(localProgram, 0755); err != nil {
										fmt.Printf(general.ErrorBaseFormat, err)
									}
								}
								text := fmt.Sprintf(general.SliceTraverse5PFormat, "==>", " ", name.(string), " ", localVersion, " -> ", remoteTag, " ", "updated")
								fmt.Printf(text)
								controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
								textLength = len(controlRegex.ReplaceAllString(text, ""))
							}
							// 生成/更新自动补全脚本
							for _, completionDir := range goCompletionDir {
								if general.FileExist(completionDir.(string)) {
									generateArgs := []string{"-c", fmt.Sprintf("%s completion zsh > %s/_%s", localProgram, completionDir.(string), name.(string))}
									if err := general.RunCommand("bash", generateArgs); err != nil {
										text := fmt.Sprintf(general.ErrorSuffixFormat, "==>", " ", acsInstallFailedMessage)
										fmt.Printf(text)
										controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
										textLength = len(controlRegex.ReplaceAllString(text, ""))
									} else {
										text := fmt.Sprintf(general.SuccessSuffixFormat, "==>", " ", acsInstallSuccessMessage)
										fmt.Printf(text)
										controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
										textLength = len(controlRegex.ReplaceAllString(text, ""))
										break
									}
								}
							}
						} else { // 压缩包校验失败
							fmt.Printf(general.ErrorSuffixFormat, "Archive file verification failed", ": ", filesInfo.ArchiveFileInfo.Name)
							continue
						}
					}
					// 分隔符和延时
					dashes := strings.Repeat("-", textLength-1)  //组装分隔符（减去行尾换行符的一个长度）
					fmt.Printf(general.LineHiddenFormat, dashes) // 美化输出
					time.Sleep(100 * time.Millisecond)           // 添加一个0.01秒的延时，使输出更加顺畅
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
					// 请求 API
					body, err := general.RequestApi(goGithubLatestSourceTagApi)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						body, err = general.RequestApi(goGiteaLatestSourceTagApi)
						if err != nil {
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
					}
					// 获取远端版本（用于 source 安装方法）
					remoteTag, err := general.GetLatestSourceTag(body)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						continue
					}
					// 获取本地版本
					localProgram := filepath.Join(installPath, name.(string)) // 本地程序路径
					nameArgs := []string{"version", "--only"}                 // 本地程序参数
					localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
					// 比较远端和本地版本
					if remoteTag == localVersion { // 版本一致，则输出无需更新信息
						text := fmt.Sprintf(general.SliceTraverse3PSuffixFormat, "==>", " ", name.(string), " ", remoteTag, " ", latestVersionMessage)
						fmt.Printf(text)
						controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
						textLength = len(controlRegex.ReplaceAllString(text, ""))
					} else { // 版本不一致，则安装或更新程序，并输出已安装/更新信息
						// 下载远端文件（如果 Temp 中已有远端文件则删除重新下载）
						goSourceTempDir := filepath.Join(installSourceTemp, name.(string))
						if general.FileExist(goSourceTempDir) {
							if err := os.RemoveAll(goSourceTempDir); err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
							}
						}
						goGithubCloneBaseUrl := fmt.Sprintf("%s/%s", goGithubUrl, goGithubUsername) // 远端仓库基础克隆地址（除仓库名）
						goGiteaCloneBaseUrl := fmt.Sprintf("%s/%s", goGiteaUrl, goGiteaUsername)    // 远端仓库基础克隆地址（除仓库名）
						fmt.Printf(general.SliceTraverse2PSuffixNoNewLineFormat, "==>", " Clone ", name.(string), " ", "from GitHub ")
						if err := cli.CloneRepoViaHTTP(installSourceTemp, goGithubCloneBaseUrl, name.(string)); err != nil {
							fmt.Printf(general.ErrorSuffixFormat, "error", " -> ", err)
							fmt.Printf(general.SliceTraverse2PSuffixNoNewLineFormat, "==>", " Clone ", name.(string), " ", "from Gitea ")
							if err := cli.CloneRepoViaHTTP(installSourceTemp, goGiteaCloneBaseUrl, name.(string)); err != nil {
								fmt.Printf(general.ErrorSuffixFormat, "error", " -> ", err)
								continue
							} else {
								fmt.Printf(general.SuccessFormat, "success")
							}
						} else {
							fmt.Printf(general.SuccessFormat, "success")
						}
						// 进到下载的远端文件目录
						if err := general.GoToDir(goSourceTempDir); err != nil {
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
						// 编译生成程序
						if general.FileExist("Makefile") { // Makefile 文件存在则使用 make 编译
							makeArgs := []string{}
							if err := general.RunCommand("make", makeArgs); err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
								continue
							}
						} else if general.FileExist("main.go") { // Makefile 文件不存在则使用 `go build` 命令编译
							buildArgs := []string{"build", "-trimpath", "-ldflags=-s -w", "-o", name.(string)}
							if err := general.RunCommand("go", buildArgs); err != nil {
								fmt.Printf(general.ErrorBaseFormat, err)
								continue
							}
						} else {
							fmt.Printf(general.ErrorBaseFormat, unableToCompileMessage)
						}
						// 检测编译生成的程序是否存在
						compileProgram := filepath.Join(installSourceTemp, name.(string), goGeneratePath, name.(string)) // 编译生成的程序
						if general.FileExist(compileProgram) {
							// 检测本地程序是否存在
							if commandErr != nil { // 不存在，安装
								if general.FileExist("Makefile") { // Makefile 文件存在则使用 `make install` 命令安装
									makeArgs := []string{"install"}
									if err := general.RunCommand("make", makeArgs); err != nil {
										fmt.Printf(general.ErrorBaseFormat, err)
										continue
									}
								} else { // Makefile 文件不存在则使用自定义函数安装
									if err := cli.InstallFile(compileProgram, localProgram, 0755); err != nil {
										fmt.Printf(general.ErrorBaseFormat, err)
										continue
									} else {
										// 为已安装的脚本设置可执行权限
										if err := os.Chmod(localProgram, 0755); err != nil {
											fmt.Printf(general.ErrorBaseFormat, err)
										}
									}
								}
								text := fmt.Sprintf(general.SliceTraverse4PFormat, "==>", " ", name.(string), " ", remoteTag, " ", "installed")
								fmt.Printf(text)
								controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
								textLength = len(controlRegex.ReplaceAllString(text, ""))
							} else { // 存在，更新
								if general.FileExist("Makefile") { // Makefile 文件存在则使用 `make install` 命令更新
									makeArgs := []string{"install"}
									if err := general.RunCommand("make", makeArgs); err != nil {
										fmt.Printf(general.ErrorBaseFormat, err)
										continue
									}
								} else { // Makefile 文件不存在则使用自定义函数更新
									if err := os.Remove(localProgram); err != nil {
										fmt.Printf(general.ErrorBaseFormat, err)
									}
									if err := cli.InstallFile(compileProgram, localProgram, 0755); err != nil {
										fmt.Printf(general.ErrorBaseFormat, err)
										continue
									} else {
										// 为已安装的脚本设置可执行权限
										if err := os.Chmod(localProgram, 0755); err != nil {
											fmt.Printf(general.ErrorBaseFormat, err)
										}
									}
								}
								text := fmt.Sprintf(general.SliceTraverse5PFormat, "==>", " ", name.(string), " ", localVersion, " -> ", remoteTag, " ", "updated")
								fmt.Printf(text)
								controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
								textLength = len(controlRegex.ReplaceAllString(text, ""))
							}
							// 生成/更新自动补全脚本
							for _, completionDir := range goCompletionDir {
								if general.FileExist(completionDir.(string)) {
									generateArgs := []string{"-c", fmt.Sprintf("%s completion zsh > %s/_%s", localProgram, completionDir.(string), name.(string))}
									if err := general.RunCommand("bash", generateArgs); err != nil {
										text := fmt.Sprintf(general.ErrorSuffixFormat, "==>", " ", acsInstallFailedMessage)
										fmt.Printf(text)
										controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
										textLength = len(controlRegex.ReplaceAllString(text, ""))
									} else {
										text := fmt.Sprintf(general.SuccessSuffixFormat, "==>", " ", acsInstallSuccessMessage)
										fmt.Printf(text)
										controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
										textLength = len(controlRegex.ReplaceAllString(text, ""))
										break
									}
								}
							}
						} else {
							text := fmt.Sprintf(general.ErrorBaseFormat, fmt.Sprintf("The source file %s does not exist", compileProgram))
							fmt.Printf(text)
							controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
							textLength = len(controlRegex.ReplaceAllString(text, ""))
						}
					}
					// 分隔符和延时
					dashes := strings.Repeat("-", textLength-1)  //组装分隔符（减去行尾换行符的一个长度）
					fmt.Printf(general.LineHiddenFormat, dashes) // 美化输出
					time.Sleep(100 * time.Millisecond)           // 添加一个0.01秒的延时，使输出更加顺畅
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
