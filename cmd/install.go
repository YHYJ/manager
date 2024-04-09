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
	Short: "Install or update programs and scripts",
	Long:  `Install or update self-developed programs and scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		cfgFile, _ := cmd.Flags().GetString("config")
		allFlag, _ := cmd.Flags().GetBool("all")
		goFlag, _ := cmd.Flags().GetBool("go")
		selfFlag, _ := cmd.Flags().GetBool("self")
		shellFlag, _ := cmd.Flags().GetBool("shell")

		var (
			noticeSlogan []string // 提示标语
		)

		// 根据参数执行操作
		if allFlag {
			goFlag, shellFlag = true, true
		}
		if selfFlag && (goFlag || shellFlag) {
			noticeSlogan = append(noticeSlogan, "参数 '--self' 只能独立使用，不能同时安装/更新其他程序")
			goFlag, shellFlag = false, false
		}

		// 读取配置文件
		configTree, err := cli.GetTomlConfig(cfgFile)
		if err != nil {
			color.Error.Println(err)
			return
		}

		// 安装/更新管理程序本身
		if selfFlag {
			color.Println()
			cli.InstallSelfProgram(configTree)
		}

		// 安装/更新基于 golang 的程序
		if goFlag {
			color.Println()
			cli.InstallGolangBasedProgram(configTree)
		}

		// 安装/更新基于 shell 的程序
		if shellFlag {
			color.Println()
			cli.InstallShellBasedProgram(configTree)
		}

		// 输出标语
		if len(noticeSlogan) > 0 {
			color.Println()
			for _, slogan := range noticeSlogan {
				color.Notice.Tips(general.PrimaryText(slogan))
			}
		}
	},
}

func init() {
	installCmd.Flags().BoolP("self", "", false, "Install or update itself (Can only be called alone)")
	installCmd.Flags().BoolP("all", "", false, "Install or update all programs and scripts")
	installCmd.Flags().BoolP("go", "", false, "Install or update programs developed based on go")
	installCmd.Flags().BoolP("shell", "", false, "Install or update shell scripts")

	installCmd.Flags().BoolP("help", "h", false, "help for install command")
	rootCmd.AddCommand(installCmd)
}
