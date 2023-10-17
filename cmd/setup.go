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
		systemcheckupdatesFlag, _ := cmd.Flags().GetBool("system-checkupdates")

		var subjectName string

		// 根据参数执行操作
		if allFlag {
			chezmoiFlag, cobraFlag, dockerFlag, frpcFlag, gitFlag, goFlag, npmFlag, pipFlag, systemcheckupdatesFlag = true, true, true, true, true, true, true, true, true
		}

		// 配置chezmoi
		if chezmoiFlag {
			subjectName = "chezmoi"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s configuration file\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.ChezmoiConfigFile)
			if err := function.WriteFile(function.ChezmoiConfigFile, function.ChezmoiConfig); err != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", err.Error())
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
			if err := function.WriteFile(function.CobraConfigFile, function.CobraConfig); err != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", err.Error())
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
			if err := function.WriteFile(function.DockerServiceConfigFile, function.DockerServiceConfig); err != nil {
				fmt.Printf("  \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n", err.Error())
			} else {
				fmt.Printf("  \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n")
			}
			// docker mirrors
			fmt.Printf(" \x1b[32m-\x1b[0m \x1b[34;1mdocker mirrors\x1b[0m\n")
			fmt.Printf("  \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s registry mirrors\x1b[0m\n", subjectName)
			fmt.Printf("  \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.DockerMirrorsConfigFile)
			if err := function.WriteFile(function.DockerMirrorsConfigFile, function.DockerMirrorsConfig); err != nil {
				fmt.Printf("  \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", err.Error())
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
			if err := function.WriteFile(function.FrpcConfigFile, function.FrpcConfig); err != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", err.Error())
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
			if err := function.WriteFile(function.GitConfigFile, function.GitConfig); err != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", err.Error())
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
			if err := function.WriteFile(function.GoConfigFile, function.GoConfig); err != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", err.Error())
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
			if err := function.WriteFile(function.NpmConfigFile, function.NpmConfig); err != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", err.Error())
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
			if err := function.WriteFile(function.PipConfigFile, function.PipConfig); err != nil {
				fmt.Printf(" \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n", err.Error())
			} else {
				fmt.Printf(" \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n\n")
			}
		}
		// 配置system-checkupdates
		if systemcheckupdatesFlag {
			// system-checkupdates timer
			subjectName = "system-checkupdates"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m \x1b[34;1m%s service\x1b[0m\n", subjectName)
			fmt.Printf("  \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s timer\x1b[0m\n", subjectName)
			fmt.Printf("  \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.SystemCheckupdatesTimerConfigFile)
			if err := function.WriteFile(function.SystemCheckupdatesTimerConfigFile, function.SystemCheckupdatesTimerConfig); err != nil {
				fmt.Printf("  \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n", err.Error())
			} else {
				fmt.Printf("  \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n")
			}
			// system-checkupdates service
			subjectName = "system-checkupdates"
			fmt.Printf("\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n", subjectName)
			fmt.Printf(" \x1b[32m-\x1b[0m \x1b[34;1m%s service\x1b[0m\n", subjectName)
			fmt.Printf("  \x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s service\x1b[0m\n", subjectName)
			fmt.Printf("  \x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n", function.SystemCheckupdatesServiceConfigFile)
			if err := function.WriteFile(function.SystemCheckupdatesServiceConfigFile, function.SystemCheckupdatesServiceConfig); err != nil {
				fmt.Printf("  \x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n", err.Error())
			} else {
				fmt.Printf("  \x1b[32m-\x1b[0m Status: \x1b[33;7mSetup completed\x1b[0m\n")
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
	setupCmd.Flags().BoolP("system-checkupdates", "", false, "Set up system-checkupdates (need to be root)")

	setupCmd.Flags().BoolP("help", "h", false, "help for setup command")
	rootCmd.AddCommand(setupCmd)
}
