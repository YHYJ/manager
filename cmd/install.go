/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 13:35:06

Description: 执行子命令 'install'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/manager/cli"
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

		// 根据参数执行操作
		if allFlag {
			goFlag, shellFlag = true, true
		}

		// 读取配置文件
		configTree, err := cli.GetTomlConfig(cfgFile)
		if err != nil {
			color.Error.Println(err)
			return
		}

		// 安装/更新基于 golang 的程序
		color.Println()
		if goFlag {
			cli.InstallGolangBasedProgram(configTree)
		}

		// 安装/更新基于 shell 的程序
		color.Println()
		if shellFlag {
			cli.InstallShellBasedProgram(configTree)
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
