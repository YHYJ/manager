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
	"github.com/yhyj/manager/general"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or update software and scripts",
	Long:  `Install or update my software and scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		configFile, _ := cmd.Flags().GetString("config")
		allFlag, _ := cmd.Flags().GetBool("all")
		goFlag, _ := cmd.Flags().GetBool("go")
		selfFlag, _ := cmd.Flags().GetBool("self")
		shellFlag, _ := cmd.Flags().GetBool("shell")

		// 根据参数执行操作
		if allFlag {
			goFlag, shellFlag = true, true
		}

		if selfFlag && (goFlag || shellFlag) {
			general.Notifier = append(general.Notifier, "'--self' cannot be mixed with other Flags")
		}

		// 读取配置文件
		configTree, err := general.GetTomlConfig(configFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Danger.Printf("Get config error (%s:%d): %s\n", fileName, lineNo+1, err)
			return
		}

		// 安装/更新管理程序本身
		if selfFlag {
			cli.InstallSelfProgram(configTree)
		}

		// 安装/更新基于 golang 的程序
		if goFlag {
			cli.InstallGolangBasedProgram(configTree)
		}

		// 安装/更新基于 shell 的程序
		if shellFlag {
			cli.InstallShellBasedProgram(configTree)
		}

		// 显示通知
		general.Notification()
	},
}

func init() {
	installCmd.Flags().BoolP("self", "", false, "Install or update itself (Can only be called alone)")
	installCmd.Flags().BoolP("all", "", false, "Install or update all software and scripts")
	installCmd.Flags().BoolP("go", "", false, "Install or update golang-based software")
	installCmd.Flags().BoolP("shell", "", false, "Install or update shell scripts")

	installCmd.Flags().BoolP("help", "h", false, "help for install command")
	rootCmd.AddCommand(installCmd)
}
