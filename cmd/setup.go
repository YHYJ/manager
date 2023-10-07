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

		var (
			subjectName string
			errInfo     error
		)

		// 根据参数执行操作
		if allFlag {
			chezmoiFlag, cobraFlag, dockerFlag, frpcFlag, gitFlag, goFlag, npmFlag, pipFlag = true, true, true, true, true, true, true, true
		}

		// 配置chezmoi
		if chezmoiFlag {
			subjectName = "chezmoi"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s configuration file\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.ChezmoiConfigFile)
			errInfo = function.WriteFile(function.ChezmoiConfigFile, function.ChezmoiConfig)
			if errInfo != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", errInfo.Error())
			} else {
				fmt.Printf(" \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
			}
		}
		// 配置cobra
		if cobraFlag {
			subjectName = "cobra-cli"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s configuration file\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.CobraConfigFile)
			errInfo = function.WriteFile(function.CobraConfigFile, function.CobraConfig)
			if errInfo != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", errInfo.Error())
			} else {
				fmt.Printf(" \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
			}
		}
		// 配置docker
		if dockerFlag {
			// docker service
			subjectName = "docker"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m \x1b[34;1mdocker service\x1b[0m\n")
			fmt.Printf("  \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s root dir\x1b[0m\n", subjectName)
			fmt.Printf("  \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.DockerServiceConfigFile)
			errInfo = function.WriteFile(function.DockerServiceConfigFile, function.DockerServiceConfig)
			if errInfo != nil {
				fmt.Printf("  \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n", errInfo.Error())
			} else {
				fmt.Printf("  \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n")
			}
			// docker mirrors
			fmt.Printf(" \x1b[32m-\x1b[0m \x1b[34;1mdocker mirrors\x1b[0m\n")
			fmt.Printf("  \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s registry mirrors\x1b[0m\n", subjectName)
			fmt.Printf("  \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.DockerMirrorsConfigFile)
			errInfo = function.WriteFile(function.DockerMirrorsConfigFile, function.DockerMirrorsConfig)
			if errInfo != nil {
				fmt.Printf("  \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", errInfo.Error())
			} else {
				fmt.Printf("  \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
			}
		}
		// 配置frpc
		if frpcFlag {
			subjectName = "frpc"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s restart timing\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.FrpcConfigFile)
			errInfo = function.WriteFile(function.FrpcConfigFile, function.FrpcConfig)
			if errInfo != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", errInfo.Error())
			} else {
				fmt.Printf(" \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
			}
		}
		// 配置git
		if gitFlag {
			subjectName = "git"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s configuration file\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.GitConfigFile)
			errInfo = function.WriteFile(function.GitConfigFile, function.GitConfig)
			if errInfo != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", errInfo.Error())
			} else {
				fmt.Printf(" \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
			}
		}
		// 配置golang
		if goFlag {
			subjectName = "go"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s environment file\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.GoConfigFile)
			errInfo = function.WriteFile(function.GoConfigFile, function.GoConfig)
			if errInfo != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", errInfo.Error())
			} else {
				fmt.Printf(" \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
			}
		}
		// 配置npm
		if npmFlag {
			subjectName = "npm"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s registry\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.NpmConfigFile)
			errInfo = function.WriteFile(function.NpmConfigFile, function.NpmConfig)
			if errInfo != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", errInfo.Error())
			} else {
				fmt.Printf(" \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
			}
		}
		// 配置pip
		if pipFlag {
			subjectName = "pip"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s mirrors\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.PipConfigFile)
			errInfo = function.WriteFile(function.PipConfigFile, function.PipConfig)
			if errInfo != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", errInfo.Error())
			} else {
				fmt.Printf(" \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
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
	setupCmd.Flags().BoolP("npm", "", false, "Set up the mirror source used by npm")
	setupCmd.Flags().BoolP("pip", "", false, "Set up the mirror source used by pip")

	setupCmd.Flags().BoolP("help", "h", false, "help for setup command")
	rootCmd.AddCommand(setupCmd)
}
