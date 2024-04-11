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
		allFlags := make(map[string]bool)
		allFlag, _ := cmd.Flags().GetBool("all")
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

		// 调用程序配置器
		cli.ProgramConfigurator(allFlags)
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
