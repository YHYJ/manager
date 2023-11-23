/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 13:35:06

Description: 程序子命令'install'时执行
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

		// 配置文件项
		var (
			installPath         string
			installSourceTemp   string
			goGeneratePath      string
			goGithubUrl         string
			goGithubApi         string
			goGithubUsername    string
			goGiteaUrl          string
			goGiteaApi          string
			goGiteaUsername     string
			goNames             []interface{}
			goCompletionDir     []interface{}
			shellGithubApi      string
			shellGithubRaw      string
			shellGithubBranch   string
			shellGithubUsername string
			shellGiteaApi       string
			shellGiteaRaw       string
			shellGiteaBranch    string
			shellGiteaUsername  string
			shellRepo           string
			shellDir            string
			shellNames          []interface{}
			httpProxy           string
			httpsProxy          string
		)

		// 输出文本
		var (
			latestVersionMessage     = "is already the latest version"                 // 已安装的程序和脚本为最新版
			unableToCompileMessage   = "Makefile or main.go file does not exist"       // 缺失编译文件无法完成编译
			acsInstallSuccessMessage = "auto-completion script installed successfully" // 自动补全脚本安装成功
			acsInstallFailedMessage  = "auto-completion script installation failed"    // 自动补全脚本安装失败
		)

		var (
			goLatestTagApiFormat             = "%s/repos/%s/%s/tags"           // 请求远端仓库最新 Tag 的 API
			shellLatestHashApiFormat         = "%s/repos/%s/%s/contents/%s/%s" // 请求远端仓库最新脚本的 Hash 值的 API
			shellGithubBaseDownloadUrlFormat = "%s/%s/%s/%s"                   // 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - Github 格式
			shellGiteaBaseDownloadUrlFormat  = "%s/%s/%s/raw/branch/%s"        // 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - Gitea 格式
		)

		// 检查配置文件是否存在
		configTree, err := cli.GetTomlConfig(cfgFile)
		if err != nil {
			fmt.Printf(general.ErrorBaseFormat, err)
			return
		} else {
			// 获取配置项
			if configTree.Has("install.path") {
				installPath = configTree.Get("install.path").(string)
			}
			if configTree.Has("install.source_temp") {
				installSourceTemp = configTree.Get("install.source_temp").(string)
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

		// 根据参数执行操作
		if allFlag {
			goFlag, shellFlag = true, true
		}

		// 安装/更新shell脚本
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
				scriptLocalPath := filepath.Join(installSourceTemp, shellRepo, name.(string))                                                              // 脚本本地存储位置
				shellGithubLatestHashApi := fmt.Sprintf(shellLatestHashApiFormat, shellGithubApi, shellGithubUsername, shellRepo, shellDir, name.(string)) // 请求远端仓库最新脚本的 Hash 值的 API
				shellGiteaLatestHashApi := fmt.Sprintf(shellLatestHashApiFormat, shellGiteaApi, shellGiteaUsername, shellRepo, shellDir, name.(string))    // 请求远端仓库最新脚本的 Hash 值的 API
				localProgram := filepath.Join(installPath, name.(string))                                                                                  // 本地程序路径
				// 请求API
				body, err := general.RequestApi(shellGithubLatestHashApi)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					body, err = general.RequestApi(shellGiteaLatestHashApi)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						continue
					}
				}
				// 获取远端脚本Hash值
				remoteHash, err := general.GetLatestSourceHash(body)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					continue
				}
				// 获取本地脚本Hash值
				gitHashObjectArgs := []string{"hash-object", localProgram} // 本地程序参数
				localHash, commandErr := general.RunCommandGetResult("git", gitHashObjectArgs)
				// 比较远端和本地脚本Hash值
				if remoteHash == localHash { // Hash值一致，则输出无需更新信息
					text := fmt.Sprintf(general.SliceTraverse2PSuffixFormat, "==>", " ", name.(string), " ", latestVersionMessage)
					fmt.Printf(text)
					controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
					textLength = len(controlRegex.ReplaceAllString(text, ""))
				} else { // Hash值不一致，则更新脚本，并输出已更新信息
					// 下载远端脚本
					shellGithubBaseDownloadUrl := fmt.Sprintf(shellGithubBaseDownloadUrlFormat, shellGithubRaw, shellGithubUsername, shellRepo, shellGithubBranch) // 脚本远端仓库基础地址
					shellGiteaBaseDownloadUrl := fmt.Sprintf(shellGiteaBaseDownloadUrlFormat, shellGiteaRaw, shellGiteaUsername, shellRepo, shellGiteaBranch)      // 脚本远端仓库基础地址
					shellUrlFile := filepath.Join(shellDir, name.(string))                                                                                         // 脚本在仓库中的路径
					fileUrl := fmt.Sprintf("%s/%s", shellGithubBaseDownloadUrl, shellUrlFile)
					_, err := cli.DownloadFile(fileUrl, scriptLocalPath)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						fileUrl := fmt.Sprintf("%s/%s", shellGiteaBaseDownloadUrl, shellUrlFile)
						_, err = cli.DownloadFile(fileUrl, scriptLocalPath)
						if err != nil {
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
				dashes := strings.Repeat("-", textLength-1)  //组装分隔符（减去行尾换行符的一个长度）
				fmt.Printf(general.LineHiddenFormat, dashes) // 美化输出
				// 添加一个0.01秒的延时，使输出更加顺畅
				time.Sleep(100 * time.Millisecond)
			}
		}

		// 安装/更新基于go开发的程序
		if goFlag {
			fmt.Printf(general.TitleH1Format, "Installing go-based programs...")
			// 设置代理
			general.SetVariable("http_proxy", httpProxy)
			general.SetVariable("https_proxy", httpsProxy)
			// 创建临时目录
			if err := general.CreateDir(installSourceTemp); err != nil {
				fmt.Printf(general.ErrorBaseFormat, err)
				return
			}
			// 遍历所有程序名
			for _, name := range goNames {
				textLength := 0                                                                                         // 输出文本的长度
				compileProgram := filepath.Join(installSourceTemp, name.(string), goGeneratePath, name.(string))        // 编译生成的最新程序
				goGithubLatestTagApi := fmt.Sprintf(goLatestTagApiFormat, goGithubApi, goGithubUsername, name.(string)) // 请求远端仓库最新 Tag 的 API
				goGiteaLatestTagApi := fmt.Sprintf(goLatestTagApiFormat, goGiteaApi, goGiteaUsername, name.(string))    // 请求远端仓库最新 Tag 的 API
				localProgram := filepath.Join(installPath, name.(string))                                               // 本地程序路径
				// 请求API
				body, err := general.RequestApi(goGithubLatestTagApi)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					body, err = general.RequestApi(goGiteaLatestTagApi)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						continue
					}
				}
				// 获取远端版本（用于Source安装）
				remoteTag, err := general.GetLatestSourceTag(body)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					continue
				}
				// 获取本地版本
				nameArgs := []string{"version", "--only"} // 本地程序参数
				localVersion, commandErr := general.RunCommandGetResult(localProgram, nameArgs)
				// 比较远端和本地版本
				if remoteTag == localVersion { // 版本一致，则输出无需更新信息
					text := fmt.Sprintf(general.SliceTraverse3PSuffixFormat, "==>", " ", name.(string), " ", remoteTag, " ", latestVersionMessage)
					fmt.Printf(text)
					controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`) // 去除控制字符，获取文本实际长度
					textLength = len(controlRegex.ReplaceAllString(text, ""))
				} else { // 版本不一致，则更新程序，并输出已更新信息
					// 下载远端文件（如果Temp中已有远端文件则删除重新下载）
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
					if general.FileExist("Makefile") { // Makefile文件存在则使用make编译
						makeArgs := []string{}
						if err := general.RunCommand("make", makeArgs); err != nil {
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
					} else if general.FileExist("main.go") { // Makefile文件不存在则使用go build编译
						buildArgs := []string{"build", "-trimpath", "-ldflags=-s -w", "-o", name.(string)}
						if err := general.RunCommand("go", buildArgs); err != nil {
							fmt.Printf(general.ErrorBaseFormat, err)
							continue
						}
					} else {
						fmt.Printf(general.ErrorBaseFormat, unableToCompileMessage)
					}
					// 检测编译生成的程序是否存在
					if general.FileExist(compileProgram) {
						// 检测本地程序是否存在
						if commandErr != nil { // 不存在，安装
							if general.FileExist("Makefile") { // Makefile文件存在则使用make install安装
								makeArgs := []string{"install"}
								if err := general.RunCommand("make", makeArgs); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
									continue
								}
							} else { // Makefile文件不存在则使用自定义函数安装
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
							if general.FileExist("Makefile") { // Makefile文件存在则使用make install更新
								makeArgs := []string{"install"}
								if err := general.RunCommand("make", makeArgs); err != nil {
									fmt.Printf(general.ErrorBaseFormat, err)
									continue
								}
							} else { // Makefile文件不存在则使用自定义函数更新
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
				dashes := strings.Repeat("-", textLength-1)  //组装分隔符（减去行尾换行符的一个长度）
				fmt.Printf(general.LineHiddenFormat, dashes) // 美化输出
				// 添加一个0.01秒的延时，使输出更加顺畅
				time.Sleep(100 * time.Millisecond)
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
