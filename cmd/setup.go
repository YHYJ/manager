/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023]-06-08 13:43:59

Description: 执行子命令 'setup'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/yhyj/manager/cli"
	"github.com/yhyj/manager/general"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up installed programs/scripts (Linux/macOS only)",
	Long:  `Set up installed self-developed programs/scripts (Linux/macOS only).`,
	Run: func(cmd *cobra.Command, args []string) {
		// 检查平台
		if general.Platform != "linux" && general.Platform != "darwin" {
			color.Printf("%s\n", general.ErrorText("Only Linux and macOS are supported"))
			return
		}

		// 解析参数
		allFlag, _ := cmd.Flags().GetBool("all")

		// 根据参数执行操作
		allFlags := make(map[string]bool)
		if allFlag {
			allFlags["chezmoiFlag"] = true
			allFlags["cobraFlag"] = true
			allFlags["dockerFlag"] = true
			allFlags["frpcFlag"] = true
			allFlags["gitFlag"] = true
			allFlags["goFlag"] = true
			allFlags["pipFlag"] = true
			allFlags["systemcheckupdatesFlag"] = true
		} else {
			allFlags["chezmoiFlag"], _ = cmd.Flags().GetBool("chezmoi")
			allFlags["cobraFlag"], _ = cmd.Flags().GetBool("cobra")
			allFlags["dockerFlag"], _ = cmd.Flags().GetBool("docker")
			allFlags["frpcFlag"], _ = cmd.Flags().GetBool("frpc")
			allFlags["gitFlag"], _ = cmd.Flags().GetBool("git")
			allFlags["goFlag"], _ = cmd.Flags().GetBool("go")
			allFlags["pipFlag"], _ = cmd.Flags().GetBool("pip")
			allFlags["systemcheckupdatesFlag"], _ = cmd.Flags().GetBool("system-checkupdates")
		}

		var (
			noticeSlogan []string // 提示标语
		)
		// 检查 allFlags 中的所有值是否都为 false
		allFalse := true
		for _, value := range allFlags {
			if value {
				allFalse = false
				break
			}
		}
		if allFalse {
			cmd.Help()
			noticeSlogan = append(noticeSlogan, "Please refer to the above help information")
		}

		// 调用程序配置器
		cli.ProgramConfigurator(allFlags)

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
	setupCmd.Flags().BoolP("all", "", false, "Set up all programs/scripts")
	setupCmd.Flags().BoolP("chezmoi", "", false, "Set up chezmoi")
	setupCmd.Flags().BoolP("cobra", "", false, "Set up cobra-cli")
	setupCmd.Flags().BoolP("docker", "", false, "Set up docker (need to be root)")
	setupCmd.Flags().BoolP("frpc", "", false, "Set up frpc restart timing (need to be root)")
	setupCmd.Flags().BoolP("git", "", false, "Set up git and generate SSH keys")
	setupCmd.Flags().BoolP("go", "", false, "Set up golang")
	setupCmd.Flags().BoolP("pip", "", false, "Set up the mirror source used by pip")
	setupCmd.Flags().BoolP("system-checkupdates", "", false, "Set up system-checkupdates (need to be root)")

	setupCmd.Flags().BoolP("help", "h", false, "help for setup command")
	rootCmd.AddCommand(setupCmd)
}
