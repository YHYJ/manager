/*
File: uninstall.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-27 12:35:06

Description: 执行子命令 'uninstall'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/manager/cli"
	"github.com/yhyj/manager/general"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall software and scripts",
	Long:  `Uninstall my software and scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		configFile, _ := cmd.Flags().GetString("config")
		allFlag, _ := cmd.Flags().GetBool("all")
		goFlag, _ := cmd.Flags().GetBool("go")
		selfFlag, _ := cmd.Flags().GetBool("self")
		shellFlag, _ := cmd.Flags().GetBool("shell")

		// 读取配置文件
		configTree, err := general.GetTomlConfig(configFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Danger.Printf("Get config error (%s:%d): %s\n", fileName, lineNo+1, err)
			return
		}

		// 根据参数执行操作
		if allFlag {
			goFlag, shellFlag = true, true
		}

		// 卸载管理程序本身
		if selfFlag {
			cli.UninstallSelf(configTree)
		}

		// 卸载基于 golang 的程序
		if goFlag {
			cli.Uninstall(configTree, "go")
		}

		// 卸载基于 shell 的程序
		if shellFlag {
			cli.Uninstall(configTree, "shell")
		}

		// 显示通知
		general.Notification()
	},
}

func init() {
	uninstallCmd.Flags().BoolP("self", "", false, "Uninstall itself")
	uninstallCmd.Flags().BoolP("all", "", false, "Uninstall all software and scripts")
	uninstallCmd.Flags().BoolP("go", "", false, "Uninstall golang-based software")
	uninstallCmd.Flags().BoolP("shell", "", false, "Uninstall shell scripts")

	uninstallCmd.Flags().BoolP("help", "h", false, "help for uninstall command")
	rootCmd.AddCommand(uninstallCmd)
}
