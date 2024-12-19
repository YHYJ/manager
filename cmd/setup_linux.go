//go:build linux

/*
File: setup_linux.go
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
		// 检查权限
		isRoot := func() bool {
			if general.GetVariable("SUDO_USER") != "" {
				return true
			}
			return false
		}()

		// 解析参数
		allFlag, _ := cmd.Flags().GetBool("all")

		// 根据参数执行操作
		allFlags := make(map[string]bool)
		if allFlag {
			if isRoot {
				allFlags["dockerFlag"] = true
				allFlags["frpcFlag"] = true
				allFlags["updatecheckerFlag"] = true
				general.Notifier = append(general.Notifier, "Please use non-root permissions to configure other")
			} else {
				allFlags["chezmoiFlag"] = true
				allFlags["cobraFlag"] = true
				allFlags["gitFlag"] = true
				allFlags["goFlag"] = true
				allFlags["pipFlag"] = true

				general.Notifier = append(general.Notifier, "Please use non-root permissions to configure other")
			}
		} else {
			allFlags["chezmoiFlag"], _ = cmd.Flags().GetBool("chezmoi")
			allFlags["cobraFlag"], _ = cmd.Flags().GetBool("cobra")
			allFlags["dockerFlag"], _ = cmd.Flags().GetBool("docker")
			allFlags["frpcFlag"], _ = cmd.Flags().GetBool("frpc")
			allFlags["gitFlag"], _ = cmd.Flags().GetBool("git")
			allFlags["goFlag"], _ = cmd.Flags().GetBool("go")
			allFlags["pipFlag"], _ = cmd.Flags().GetBool("pip")
			allFlags["updatecheckerFlag"], _ = cmd.Flags().GetBool("update-checker")
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
	setupCmd.Flags().Bool("docker", false, "Set up docker (need to be root)")
	setupCmd.Flags().Bool("frpc", false, "Set up frpc restart timing (need to be root)")
	setupCmd.Flags().Bool("git", false, "Set up git and generate SSH keys")
	setupCmd.Flags().Bool("go", false, "Set up golang")
	setupCmd.Flags().Bool("pip", false, "Set up the mirror source used by pip")
	setupCmd.Flags().Bool("update-checker", false, "Set up update-checker (need to be root)")

	setupCmd.Flags().BoolP("help", "h", false, "help for setup command")
	rootCmd.AddCommand(setupCmd)
}
