/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023]-06-08 13:43:59

Description: 程序子命令'setup'时执行
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yhyj/manager/function"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up installed programs/scripts",
	Long:  `Set up installed self-developed programs/scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		allFlag, _ := cmd.Flags().GetBool("all")
		pipFlag, _ := cmd.Flags().GetBool("pip")
		npmFlag, _ := cmd.Flags().GetBool("npm")
		dockerFlag, _ := cmd.Flags().GetBool("docker")
		gitFlag, _ := cmd.Flags().GetBool("git")

		var errSubject string
		var errInfo error
		var errReport string

		// 根据参数执行操作
		if allFlag {
			pipFlag, npmFlag, dockerFlag, gitFlag = true, true, true, true
		}
		if pipFlag {
			// 配置pip
			errSubject = "pip"
			errInfo = function.WriteFile(function.PipConfigFile, function.PipConfig)
			if errInfo != nil {
				errReport = errReport + errSubject + ": " + errInfo.Error() + "\n"
			}
		}
		if npmFlag {
			// 配置npm
			errSubject = "npm"
			errInfo = function.WriteFile(function.PipConfigFile, function.PipConfig)
			if errInfo != nil {
				errReport = errReport + errSubject + ": " + errInfo.Error() + "\n"
			}
		}
		if dockerFlag {
			// 配置docker
			errSubject = "docker"
			errInfo = function.WriteFile(function.PipConfigFile, function.PipConfig)
			if errInfo != nil {
				errReport = errReport + errSubject + ": " + errInfo.Error() + "\n"
			}
		}
		if gitFlag {
			// 配置git
			errSubject = "git"
			errInfo = function.WriteFile(function.PipConfigFile, function.PipConfig)
			if errInfo != nil {
				errReport = errReport + errSubject + ": " + errInfo.Error() + "\n"
			}
		}

		if errReport != "" {
			fmt.Printf("\x1b[36;1m%s\x1b[0m\n", errReport)
		}
	},
}

func init() {
	setupCmd.Flags().BoolP("all", "", false, "Set up all programs/scripts")
	setupCmd.Flags().BoolP("pip", "", false, "Set up the mirror source used by pip")
	setupCmd.Flags().BoolP("npm", "", false, "Set up the mirror source used by npm")
	setupCmd.Flags().BoolP("docker", "", false, "Set up Docker Root Directory")
	setupCmd.Flags().BoolP("git", "", false, "Set up git and generate SSH keys")

	setupCmd.Flags().BoolP("help", "h", false, "help for setup")
	rootCmd.AddCommand(setupCmd)
}
