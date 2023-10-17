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

		// 根据参数执行操作
		if allFlag {
			chezmoiFlag, cobraFlag, dockerFlag, frpcFlag, gitFlag, goFlag, npmFlag, pipFlag, systemcheckupdatesFlag = true, true, true, true, true, true, true, true, true
		}

		// 预定义变量
		var (
			subjectName      string
			descriptorText   string
			subjectMinorName string
		)

		// 定义输出格式
		subjectNameFormat := "\x1b[32;1m==>\x1b[0m \x1b[34m%s\x1b[0m\n"
		subjectMinorNameFormat := "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m \x1b[34;1m%s\x1b[0m\n"
		descriptorFormat := "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m Descriptor: \x1b[33mSet up %s %s\x1b[0m\n"
		configFileFormat := "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m Configuration file: \x1b[33m%s\x1b[0m\n"
		errorFormat := "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n\n"
		successFormat := "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m Status: \x1b[33;7m%s\x1b[0m\n\n"

		// 配置chezmoi
		if chezmoiFlag {
			subjectName = "chezmoi"
			descriptorText = "configuration file"
			fmt.Printf(subjectNameFormat, subjectName)
			fmt.Printf(descriptorFormat, 1, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 1, function.ChezmoiConfigFile)
			if err := function.WriteFile(function.ChezmoiConfigFile, function.ChezmoiConfig); err != nil {
				fmt.Printf(errorFormat, 1, err.Error())
			} else {
				fmt.Printf(successFormat, 1, "Setup completed")
			}
		}
		// 配置cobra
		if cobraFlag {
			subjectName = "cobra-cli"
			descriptorText = "configuration file"
			fmt.Printf(subjectNameFormat, subjectName)
			fmt.Printf(descriptorFormat, 1, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 1, function.CobraConfigFile)
			if err := function.WriteFile(function.CobraConfigFile, function.CobraConfig); err != nil {
				fmt.Printf(errorFormat, 1, err.Error())
			} else {
				fmt.Printf(successFormat, 1, "Setup completed")
			}
		}
		// 配置docker
		if dockerFlag {
			subjectName = "docker"
			fmt.Printf(subjectNameFormat, subjectName)
			// docker service
			subjectMinorName = "docker service"
			descriptorText = "root directory"
			fmt.Printf(subjectMinorNameFormat, 1, subjectMinorName)
			fmt.Printf(descriptorFormat, 2, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 2, function.DockerServiceConfigFile)
			if err := function.WriteFile(function.DockerServiceConfigFile, function.DockerServiceConfig); err != nil {
				errorFormat = "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n"
				fmt.Printf(errorFormat, 2, err.Error())
			} else {
				successFormat := "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m Status: \x1b[33;7m%s\x1b[0m\n"
				fmt.Printf(successFormat, 2, "Setup completed")
			}
			// docker mirrors
			subjectMinorName = "docker mirrors"
			descriptorText = "registry mirrors"
			fmt.Printf(subjectMinorNameFormat, 1, subjectMinorName)
			fmt.Printf(descriptorFormat, 2, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 2, function.DockerMirrorsConfigFile)
			if err := function.WriteFile(function.DockerMirrorsConfigFile, function.DockerMirrorsConfig); err != nil {
				fmt.Printf(errorFormat, 2, err.Error())
			} else {
				fmt.Printf(successFormat, 2, "Setup completed")
			}
		}
		// 配置frpc
		if frpcFlag {
			subjectName = "frpc"
			descriptorText = "restart timing"
			fmt.Printf(subjectNameFormat, subjectName)
			fmt.Printf(descriptorFormat, 1, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 1, function.FrpcConfigFile)
			if err := function.WriteFile(function.FrpcConfigFile, function.FrpcConfig); err != nil {
				fmt.Printf(errorFormat, 1, err.Error())
			} else {
				fmt.Printf(successFormat, 1, "Setup completed")
			}
		}
		// 配置git
		if gitFlag {
			subjectName = "git"
			descriptorText = "configuration file"
			fmt.Printf(subjectNameFormat, subjectName)
			fmt.Printf(descriptorFormat, 1, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 1, function.GitConfigFile)
			if err := function.WriteFile(function.GitConfigFile, function.GitConfig); err != nil {
				fmt.Printf(errorFormat, 1, err.Error())
			} else {
				fmt.Printf(successFormat, 1, "Setup completed")
			}
		}
		// 配置golang
		if goFlag {
			subjectName = "go"
			descriptorText = "environment file"
			fmt.Printf(subjectNameFormat, subjectName)
			fmt.Printf(descriptorFormat, 1, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 1, function.GoConfigFile)
			if err := function.WriteFile(function.GoConfigFile, function.GoConfig); err != nil {
				fmt.Printf(errorFormat, 1, err.Error())
			} else {
				fmt.Printf(successFormat, 1, "Setup completed")
			}
		}
		// 配置npm
		if npmFlag {
			subjectName = "npm"
			descriptorText = "registry"
			fmt.Printf(subjectNameFormat, subjectName)
			fmt.Printf(descriptorFormat, 1, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 1, function.NpmConfigFile)
			if err := function.WriteFile(function.NpmConfigFile, function.NpmConfig); err != nil {
				fmt.Printf(errorFormat, 1, err.Error())
			} else {
				fmt.Printf(successFormat, 1, "Setup completed")
			}
		}
		// 配置pip
		if pipFlag {
			subjectName = "pip"
			descriptorText = "mirrors"
			fmt.Printf(subjectNameFormat, subjectName)
			fmt.Printf(descriptorFormat, 1, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 1, function.PipConfigFile)
			if err := function.WriteFile(function.PipConfigFile, function.PipConfig); err != nil {
				fmt.Printf(errorFormat, 1, err.Error())
			} else {
				fmt.Printf(successFormat, 1, "Setup completed")
			}
		}
		// 配置system-checkupdates
		if systemcheckupdatesFlag {
			subjectName = "system-checkupdates"
			fmt.Printf(subjectNameFormat, subjectName)
			// system-checkupdates timer
			subjectMinorName = "system-checkupdates timer"
			descriptorText = "timer"
			fmt.Printf(subjectMinorNameFormat, 1, subjectMinorName)
			fmt.Printf(descriptorFormat, 2, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 2, function.SystemCheckupdatesTimerConfigFile)
			if err := function.WriteFile(function.SystemCheckupdatesTimerConfigFile, function.SystemCheckupdatesTimerConfig); err != nil {
				errorFormat = "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m Error: \x1b[31m%s\x1b[0m\n"
				fmt.Printf(errorFormat, 2, err.Error())
			} else {
				successFormat := "\x1b[%dC\x1b[0m\x1b[32m-\x1b[0m Status: \x1b[33;7m%s\x1b[0m\n"
				fmt.Printf(successFormat, 2, "Setup completed")
			}
			// system-checkupdates service
			subjectMinorName = "system-checkupdates service"
			descriptorText = "service"
			fmt.Printf(subjectMinorNameFormat, 1, subjectMinorName)
			fmt.Printf(descriptorFormat, 2, subjectName, descriptorText)
			fmt.Printf(configFileFormat, 2, function.SystemCheckupdatesServiceConfigFile)
			if err := function.WriteFile(function.SystemCheckupdatesServiceConfigFile, function.SystemCheckupdatesServiceConfig); err != nil {
				fmt.Printf(errorFormat, 2, err.Error())
			} else {
				fmt.Printf(successFormat, 2, "Setup completed")
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
