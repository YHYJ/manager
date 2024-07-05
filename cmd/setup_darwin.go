//go:build darwin

/*
File: setup_darwin.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023]-06-08 13:43:59

Description: 执行子命令 'setup'
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yhyj/manager/cli"
	"github.com/yhyj/manager/general"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up installed programs/scripts",
	Long:  `Set up installed self-developed programs/scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		allFlag, _ := cmd.Flags().GetBool("all")

		// 根据参数执行操作
		allFlags := make(map[string]bool)
		if allFlag {
			allFlags["chezmoiFlag"] = true
			allFlags["cobraFlag"] = true
			allFlags["gitFlag"] = true
			allFlags["goFlag"] = true
			allFlags["pipFlag"] = true
		} else {
			allFlags["chezmoiFlag"], _ = cmd.Flags().GetBool("chezmoi")
			allFlags["cobraFlag"], _ = cmd.Flags().GetBool("cobra")
			allFlags["gitFlag"], _ = cmd.Flags().GetBool("git")
			allFlags["goFlag"], _ = cmd.Flags().GetBool("go")
			allFlags["pipFlag"], _ = cmd.Flags().GetBool("pip")
		}

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
			general.Notifier = append(general.Notifier, "Please refer to the above help information")
		}

		// 调用程序配置器
		cli.ProgramConfigurator(allFlags)

		// 显示通知
		general.Notification()
	},
}

func init() {
	setupCmd.Flags().Bool("all", false, "Set up all programs/scripts")
	setupCmd.Flags().Bool("chezmoi", false, "Set up chezmoi")
	setupCmd.Flags().Bool("cobra", false, "Set up cobra-cli")
	setupCmd.Flags().Bool("git", false, "Set up git and generate SSH keys")
	setupCmd.Flags().Bool("go", false, "Set up golang")
	setupCmd.Flags().Bool("pip", false, "Set up the mirror source used by pip")

	setupCmd.Flags().BoolP("help", "h", false, "help for setup command")
	rootCmd.AddCommand(setupCmd)
}
