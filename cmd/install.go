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
			installPath     string
			installTemp     string
			goSource        string
			goNames         []interface{}
			goCompletionDir string
			shellSource     string
			shellRepo       string
			shellDir        string
			shellNames      []interface{}
			httpProxy       string
			httpsProxy      string
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
			if configTree.Has("install.go.source") {
				goSource = configTree.Get("install.go.source").(string)
			}
			if configTree.Has("install.go.names") {
				goNames = configTree.Get("install.go.names").([]interface{})
			}
			if configTree.Has("install.go.completion_dir") {
				goCompletionDir = configTree.Get("install.go.completion_dir").(string)
			}
			if configTree.Has("install.shell.source") {
				shellSource = configTree.Get("install.shell.source").(string)
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
				err := function.CreateDir(installTemp)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				// 下载源文件（如果Temp中已有源文件则删除重新下载）
				if function.FileExist(installTemp + "/" + shellRepo) {
					if err := os.RemoveAll(installTemp + "/" + shellRepo); err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
						return
					}
				}
				function.CloneRepoViaHTTP(installTemp, shellSource, shellRepo)
				// 进到源文件目录
				err = function.GoToDir(installTemp + "/" + shellRepo)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				// 遍历所有脚本名
				for _, name := range shellNames {
					// 组装文件名变量
					tempAreaFile := installTemp + "/" + shellRepo + shellDir + "/" + name.(string) // 最新文件
					pathAreaFile := installPath + "/" + name.(string)                              // 已安装文件
					// 检测源文件是否存在
					if function.FileExist(tempAreaFile) {
						// 检测目标文件是否存在
						if !function.FileExist(pathAreaFile) { // 不存在，安装
							err := function.InstallFile(tempAreaFile, pathAreaFile)
							if err != nil {
								fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
							} else {
								fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m \x1b[35;1minstallation\x1b[0m complete\n", name.(string))
							}
						} else { // 存在，更新
							// 判断已安装的程序和要安装的文件是否一样
							equal, err := function.CompareFile(tempAreaFile, pathAreaFile)
							if err != nil {
								fmt.Printf("\x1b[31mCompare file error: %s\x1b[0m\n", err)
								return
							}
							if equal {
								// 一样，则输出无需更新信息
								fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m is already the latest version\n", name.(string))
							} else {
								// 不一样，则更新程序，并输出已更新信息
								fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m \x1b[35;1mupdate\x1b[0m complete\n", name.(string))
								if err := os.Remove(pathAreaFile); err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
									return
								}
								err := function.InstallFile(tempAreaFile, pathAreaFile)
								if err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
								}
							}
						}
					} else {
						fmt.Printf("\x1b[31mThe source file %s does not exist\x1b[0m\n", tempAreaFile)
					}
				}
			}
			// 安装/更新基于go开发的程序
			if goFlag {
				fmt.Printf("\x1b[36;3m%s\x1b[0m\n", "Installing go-based programs...")
				// 设置代理
				function.SetVariable("http_proxy", httpProxy)
				function.SetVariable("https_proxy", httpsProxy)
				// 创建临时目录
				err := function.CreateDir(installTemp)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				// 遍历所有程序名
				for _, name := range goNames {
					// 组装文件名变量
					tempAreaFile := installTemp + "/" + name.(string) + "/" + name.(string) // 最新文件
					pathAreaFile := installPath + "/" + name.(string)                       // 已安装文件
					// 下载源文件（如果Temp中已有源文件则删除重新下载）
					if function.FileExist(installTemp + "/" + name.(string)) {
						if err := os.RemoveAll(installTemp + "/" + name.(string)); err != nil {
							fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
							return
						}
					}
					function.CloneRepoViaHTTP(installTemp, goSource, name.(string))
					// 进到源文件目录
					err := function.GoToDir(installTemp + "/" + name.(string))
					if err != nil {
						fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
						return
					}
					// 编译生成二进制文件
					buildArgs := []string{"build", "-trimpath", "-o", name.(string)}
					function.RunCommandGetFlag("go", buildArgs)
					// 检测源文件是否存在
					if function.FileExist(tempAreaFile) {
						// 检测目标文件是否存在
						if !function.FileExist(pathAreaFile) { // 不存在，安装
							err := function.InstallFile(tempAreaFile, pathAreaFile)
							if err != nil {
								fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
							} else {
								fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m \x1b[35;1minstallation\x1b[0m complete\n", name.(string))
							}
						} else { // 存在，更新
							// 判断已安装的程序和要安装的文件是否一样
							equal, err := function.CompareFile(tempAreaFile, pathAreaFile)
							if err != nil {
								fmt.Printf("\x1b[31mCompare file error: %s\x1b[0m\n", err)
								return
							}
							if equal {
								// 一样，则输出无需更新信息
								fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m is already the latest version\n", name.(string))
							} else {
								// 不一样，则更新程序，并输出已更新信息
								fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m \x1b[35;1mupdate\x1b[0m complete\n", name.(string))
								if err := os.Remove(pathAreaFile); err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
									return
								}
								err := function.InstallFile(tempAreaFile, pathAreaFile)
								if err != nil {
									fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
								}
							}
						}
					} else {
						fmt.Printf("\x1b[31mThe source file %s does not exist\x1b[0m\n", tempAreaFile)
					}
					// 生成/更新自动补全脚本
					generateArgs := []string{"-c", fmt.Sprintf("%s completion zsh > %s", pathAreaFile, goCompletionDir+"/"+"_"+name.(string))}
					flag := function.RunCommandGetFlag("bash", generateArgs)
					if flag {
						fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m auto-completion script installed successfully\n\n", name.(string))
					} else {
						fmt.Printf("\x1b[31m==>\x1b[0m \x1b[34m%s\x1b[0m auto-completion script installation failed\n\n", name.(string))
					}
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
