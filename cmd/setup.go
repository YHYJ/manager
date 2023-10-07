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
		chezmoiFlag, _ := cmd.Flags().GetBool("chezmoi")
		cobraFlag, _ := cmd.Flags().GetBool("cobra")
		dockerFlag, _ := cmd.Flags().GetBool("docker")
		frpcFlag, _ := cmd.Flags().GetBool("frpc")
		gitFlag, _ := cmd.Flags().GetBool("git")
		goFlag, _ := cmd.Flags().GetBool("go")
		npmFlag, _ := cmd.Flags().GetBool("npm")
		pipFlag, _ := cmd.Flags().GetBool("pip")

		// 接收错误信息的变量
		var (
		errSubject string
		errInfo error
		errReport string
		)

		// 根据参数执行操作
		if allFlag {
			chezmoiFlag, cobraFlag, dockerFlag, frpcFlag, gitFlag, goFlag, npmFlag, pipFlag = true, true, true, true, true, true, true, true
		}

		// 配置chezmoi
		if chezmoiFlag {
			errSubject = "chezmoi"
			errInfo = function.WriteFile(function.ChezmoiConfigFile, function.ChezmoiConfig)
			if errInfo != nil {
				errReport = fmt.Sprintf("%s: %s\n", errSubject, errInfo.Error())
				fmt.Printf("\x1b[31m%s\x1b[0m\n", errReport)
			}
		}
		// 配置cobra
		if cobraFlag {
			errSubject = "cobra"
			errInfo = function.WriteFile(function.CobraConfigFile, function.CobraConfig)
			if errInfo != nil {
				errReport = fmt.Sprintf("%s: %s\n", errSubject, errInfo.Error())
				fmt.Printf("\x1b[31m%s\x1b[0m\n", errReport)
			}
		}
		// 配置docker
		if dockerFlag {
			errSubject = "docker"
			errInfo = function.WriteFile(function.DockerConfigFile, function.DockerConfig)
			if errInfo != nil {
				errReport = fmt.Sprintf("%s: %s\n", errSubject, errInfo.Error())
				fmt.Printf("\x1b[31m%s\x1b[0m\n", errReport)
			}
		}
		// 配置frpc
		if frpcFlag {
			errSubject = "frpc"
			errInfo = function.WriteFile(function.FrpcConfigFile, function.FrpcConfig)
			if errInfo != nil {
				errReport = fmt.Sprintf("%s: %s\n", errSubject, errInfo.Error())
				fmt.Printf("\x1b[31m%s\x1b[0m\n", errReport)
			}
		}
		// 配置git
		if gitFlag {
			errSubject = "git"
			errInfo = function.WriteFile(function.GitConfigFile, function.GitConfig)
			if errInfo != nil {
				errReport = fmt.Sprintf("%s: %s\n", errSubject, errInfo.Error())
				fmt.Printf("\x1b[31m%s\x1b[0m\n", errReport)
			}
		}
		// 配置golang
		if goFlag {
			errSubject = "go"
			errInfo = function.WriteFile(function.GoConfigFile, function.GoConfig)
			if errInfo != nil {
				errReport = fmt.Sprintf("%s: %s\n", errSubject, errInfo.Error())
				fmt.Printf("\x1b[31m%s\x1b[0m\n", errReport)
			}
		}
		// 配置npm
		if npmFlag {
			errSubject = "npm"
			errInfo = function.WriteFile(function.NpmConfigFile, function.NpmConfig)
			if errInfo != nil {
				errReport = fmt.Sprintf("%s: %s\n", errSubject, errInfo.Error())
				fmt.Printf("\x1b[31m%s\x1b[0m\n", errReport)
			}
		}
		// 配置pip
		if pipFlag {
			errSubject = "pip"
			errInfo = function.WriteFile(function.PipConfigFile, function.PipConfig)
			if errInfo != nil {
				errReport = fmt.Sprintf("%s: %s\n", errSubject, errInfo.Error())
				fmt.Printf("\x1b[31m%s\x1b[0m\n", errReport)
			}
		}
	},
}

func init() {
	setupCmd.Flags().BoolP("all", "", false, "Set up all programs/scripts")
	setupCmd.Flags().BoolP("chezmoi", "", false, "Set up chezmoi")
	setupCmd.Flags().BoolP("cobra", "", false, "Set up cobra-cli")
	setupCmd.Flags().BoolP("docker", "", false, "Set up Docker Root Directory (need to be root)")
	setupCmd.Flags().BoolP("frpc", "", false, "Set up frpc restart timing (need to be root)")
	setupCmd.Flags().BoolP("git", "", false, "Set up git and generate SSH keys")
	setupCmd.Flags().BoolP("go", "", false, "Set up golang")
	setupCmd.Flags().BoolP("npm", "", false, "Set up the mirror source used by npm")
	setupCmd.Flags().BoolP("pip", "", false, "Set up the mirror source used by pip")

	setupCmd.Flags().BoolP("help", "h", false, "help for setup command")
	rootCmd.AddCommand(setupCmd)
}
