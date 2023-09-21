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

	"github.com/spf13/cobra"
	"github.com/yhyj/manager/function"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or update programs and scripts",
	Long:  `Install or update self-developed programs and scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取配置文件路径
		cfgFile, _ := cmd.Flags().GetString("config")
		// 解析参数
		allFlag, _ := cmd.Flags().GetBool("all")
		goFlag, _ := cmd.Flags().GetBool("go")
		shellFlag, _ := cmd.Flags().GetBool("shell")

		var (
			installPath         string
			installTemp         string
			goSourceUrl         string
			goSourceApi         string
			goSourceUsername    string
			goNames             []interface{}
			goCompletionDir     string
			shellSourceUrl      string
			shellSourceApi      string
			shellSourceBranch   string
			shellSourceUsername string
			shellRepo           string
			shellDir            string
			shellNames          []interface{}
			httpProxy           string
			httpsProxy          string
		)
		// 检查配置文件是否存在
		configTree, err := function.GetTomlConfig(cfgFile)
		if err != nil {
			fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
		} else {
			// 获取配置项
			if configTree.Has("install.path") {
				installPath = configTree.Get("install.path").(string)
			}
			if configTree.Has("install.temp") {
				installTemp = configTree.Get("install.temp").(string)
			}
			if configTree.Has("install.go.source_url") {
				goSourceUrl = configTree.Get("install.go.source_url").(string)
			}
			if configTree.Has("install.go.source_api") {
				goSourceApi = configTree.Get("install.go.source_api").(string)
			}
			if configTree.Has("install.go.source_username") {
				goSourceUsername = configTree.Get("install.go.source_username").(string)
			}
			if configTree.Has("install.go.names") {
				goNames = configTree.Get("install.go.names").([]interface{})
			}
			if configTree.Has("install.go.completion_dir") {
				goCompletionDir = configTree.Get("install.go.completion_dir").(string)
			}
			if configTree.Has("install.shell.source_url") {
				shellSourceUrl = configTree.Get("install.shell.source_url").(string)
			}
			if configTree.Has("install.shell.source_api") {
				shellSourceApi = configTree.Get("install.shell.source_api").(string)
			}
			if configTree.Has("install.shell.source_branch") {
				shellSourceBranch = configTree.Get("install.shell.source_branch").(string)
			}
			if configTree.Has("install.shell.source_username") {
				shellSourceUsername = configTree.Get("install.shell.source_username").(string)
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

			// 根据参数执行操作
			if allFlag {
				goFlag, shellFlag = true, true
			}
			// 安装/更新shell脚本
			if shellFlag {
				fmt.Printf("\x1b[36;3m%s\x1b[0m\n", "Installing shell-based programs...")
				// 设置代理
				function.SetVariable("http_proxy", httpProxy)
				function.SetVariable("https_proxy", httpsProxy)
				// 创建临时目录
				if err := function.CreateDir(installTemp); err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				// 遍历所有脚本名
				for _, name := range shellNames {
					// 组装变量
					compileProgram := fmt.Sprintf("%s/%s/%s", installTemp, shellRepo, name.(string))                                                           // 从远端下载的最新脚本
					shellSourceApiUrl := fmt.Sprintf("%s/repos/%s/%s/contents/%s/%s", shellSourceApi, shellSourceUsername, shellRepo, shellDir, name.(string)) // API URL
					localProgram := fmt.Sprintf("%s/%s", installPath, name.(string))                                                                           // 本地程序路径
					gitHashObjectArgs := []string{"hash-object", localProgram}                                                                                 // 本地程序参数
					// 请求API
					body, err := function.RequestApi(shellSourceApiUrl)
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
						return
					}
					// 获取远端脚本Hash值
					remoteHash, err := function.ParseApiResponse(body)
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
						return
					}
					// 获取本地脚本Hash值
					localHash, commandErr := function.RunCommandGetResult("git", gitHashObjectArgs)
					// 比较远端和本地脚本Hash值
					if remoteHash == localHash { // Hash值一致，则输出无需更新信息
						fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m is already the latest version\n", name.(string))
					} else { // Hash值不一致，则更新脚本，并输出已更新信息
						// 下载远端脚本
						shellSourceTempDir := fmt.Sprintf("%s/%s", installTemp, shellRepo)
						shellSource := fmt.Sprintf("%s/%s/%s/raw/branch/%s", shellSourceUrl, shellSourceUsername, shellRepo, shellSourceBranch)
						shellUrlFile := fmt.Sprintf("%s/%s", shellDir, name.(string))
						shellOutputFile := fmt.Sprintf("%s/%s", shellSourceTempDir, name.(string))
						function.DownloadFile(shellSource, shellUrlFile, shellOutputFile)
						// 检测脚本文件是否存在
						if function.FileExist(compileProgram) {
							// 检测本地程序是否存在
							if commandErr != nil { // 不存在，安装
								if err := function.InstallFile(compileProgram, localProgram); err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
								} else {
									// 为已安装的脚本设置可执行权限
									if err = os.Chmod(localProgram, 0755); err != nil {
										fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
									}
									fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m \x1b[35;1minstallation\x1b[0m complete\n", name.(string))
								}
							} else { // 存在，更新
								if err := os.Remove(localProgram); err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
									return
								}
								if err := function.InstallFile(compileProgram, localProgram); err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
								} else {
									// 为已更新的脚本设置可执行权限
									if err = os.Chmod(localProgram, 0755); err != nil {
										fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
									}
									fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m \x1b[35;1mupdate\x1b[0m complete\n", name.(string))
								}
							}
						} else {
							fmt.Printf("\x1b[31mThe source file %s does not exist\x1b[0m\n", compileProgram)
						}
					}
					fmt.Println()
				}
			}
			// 安装/更新基于go开发的程序
			if goFlag {
				fmt.Printf("\x1b[36;3m%s\x1b[0m\n", "Installing go-based programs...")
				// 设置代理
				function.SetVariable("http_proxy", httpProxy)
				function.SetVariable("https_proxy", httpsProxy)
				// 创建临时目录
				if err := function.CreateDir(installTemp); err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				// 遍历所有程序名
				for _, name := range goNames {
					// 组装变量
					compileProgram := fmt.Sprintf("%s/%s/%s", installTemp, name.(string), name.(string))               // 编译生成的最新程序
					goSourceApiUrl := fmt.Sprintf("%s/repos/%s/%s/tags", goSourceApi, goSourceUsername, name.(string)) // API URL
					localProgram := fmt.Sprintf("%s/%s", installPath, name.(string))                                   // 本地程序路径
					nameArgs := []string{"version", "--only"}                                                          // 本地程序参数
					// 请求API
					body, err := function.RequestApi(goSourceApiUrl)
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
						return
					}
					// 获取远端版本
					remoteVersion, err := function.ParseApiResponse(body)
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
						return
					}
					// 获取本地版本
					localVersion, commandErr := function.RunCommandGetResult(localProgram, nameArgs)
					// 比较远端和本地版本
					if remoteVersion == localVersion { // 版本一致，则输出无需更新信息
						fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m is already the latest version\n", name.(string))
					} else { // 版本不一致，则更新程序，并输出已更新信息
						// 下载远端文件（如果Temp中已有远端文件则删除重新下载）
						goSourceTempDir := fmt.Sprintf("%s/%s", installTemp, name.(string))
						if function.FileExist(goSourceTempDir) {
							if err := os.RemoveAll(goSourceTempDir); err != nil {
								fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
								return
							}
						}
						goSource := fmt.Sprintf("%s/%s", goSourceUrl, goSourceUsername)
						function.CloneRepoViaHTTP(installTemp, goSource, name.(string))
						// 进到下载的远端文件目录
						if err = function.GoToDir(goSourceTempDir); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
							return
						}
						// 编译生成程序
						buildArgs := []string{"build", "-trimpath", "-o", name.(string)}
						function.RunCommandGetFlag("go", buildArgs)
						// 检测编译生成的程序是否存在
						if function.FileExist(compileProgram) {
							// 检测本地程序是否存在
							if commandErr != nil { // 不存在，安装
								if err := function.InstallFile(compileProgram, localProgram); err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
								} else {
									fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m \x1b[35;1minstallation\x1b[0m complete\n", name.(string))
								}
							} else { // 存在，更新
								if err := os.Remove(localProgram); err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
									return
								}
								if err := function.InstallFile(compileProgram, localProgram); err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
								} else {
									fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m \x1b[35;1mupdate\x1b[0m complete\n", name.(string))
								}
							}
						} else {
							fmt.Printf("\x1b[31mThe source file %s does not exist\x1b[0m\n", compileProgram)
						}
						// 生成/更新自动补全脚本
						copmleteFile := fmt.Sprintf("%s/_%s", goCompletionDir, name.(string))
						generateArgs := []string{"-c", fmt.Sprintf("%s completion zsh > %s", localProgram, copmleteFile)}
						flag := function.RunCommandGetFlag("bash", generateArgs)
						if flag {
							fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m auto-completion script installed successfully\n", name.(string))
						} else {
							fmt.Printf("\x1b[31m==>\x1b[0m \x1b[34m%s\x1b[0m auto-completion script installation failed\n", name.(string))
						}
					}
					fmt.Println()
				}
			}
		}
	},
}

func init() {
	installCmd.Flags().BoolP("all", "", false, "Install or update all programs and scripts")
	installCmd.Flags().BoolP("go", "", false, "Install or update programs developed based on go")
	installCmd.Flags().BoolP("shell", "", false, "Install or update shell scripts")

	installCmd.Flags().BoolP("help", "h", false, "help for install")
	rootCmd.AddCommand(installCmd)
}
