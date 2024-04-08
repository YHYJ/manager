/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023]-06-08 13:43:59

Description: 执行子命令 'setup'
*/

package cmd

import (
	"fmt"

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
		chezmoiFlag, _ := cmd.Flags().GetBool("chezmoi")
		cobraFlag, _ := cmd.Flags().GetBool("cobra")
		dockerFlag, _ := cmd.Flags().GetBool("docker")
		frpcFlag, _ := cmd.Flags().GetBool("frpc")
		gitFlag, _ := cmd.Flags().GetBool("git")
		goFlag, _ := cmd.Flags().GetBool("go")
		pipFlag, _ := cmd.Flags().GetBool("pip")
		systemcheckupdatesFlag, _ := cmd.Flags().GetBool("system-checkupdates")

		// 根据参数执行操作
		if allFlag {
			chezmoiFlag, cobraFlag, dockerFlag, frpcFlag, gitFlag, goFlag, pipFlag, systemcheckupdatesFlag = true, true, true, true, true, true, true, true
		}

		// 预定义变量
		var (
			subjectName      string
			descriptorText   string
			subjectMinorName string
		)

		// 定义输出格式
		subjectMinorNameFormat := "%*s%s %s\n"
		descriptorFormat := "%*s%s %s: %s %s %s\n"
		configFileFormat := "%*s%s %s: %s\n"
		errorFormat := "%*s%s %s: %s\n\n"
		successFormat := "%*s%s %s: %s\n\n"

		// 配置 chezmoi
		if chezmoiFlag {
			subjectName = "chezmoi"
			descriptorText = "configuration file"
			color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlue(subjectName))
			color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.ChezmoiConfigFile))
			if err := general.WriteFile(cli.ChezmoiConfigFile, cli.ChezmoiConfig); err != nil {
				color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				fmt.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
		}
		// 配置 cobra
		if cobraFlag {
			subjectName = "cobra-cli"
			descriptorText = "configuration file"
			color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlue(subjectName))
			color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.CobraConfigFile))
			if err := general.WriteFile(cli.CobraConfigFile, cli.CobraConfig); err != nil {
				color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				fmt.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
		}
		// 配置 docker
		if dockerFlag {
			subjectName = "docker"
			color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlue(subjectName))
			// docker service
			subjectMinorName = "docker service"
			descriptorText = "root directory"
			color.Printf(subjectMinorNameFormat, 1, " ", general.SuccessText("-"), general.FgBlue(subjectMinorName))
			color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 2, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.DockerServiceConfigFile))
			if err := general.WriteFile(cli.DockerServiceConfigFile, cli.DockerServiceConfig); err != nil {
				errorFormat = "%*s%s %s: %s\n"
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				successFormat = "%*s%s %s: %s\n"
				fmt.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
			// docker mirrors
			subjectMinorName = "docker mirrors"
			descriptorText = "registry mirrors"
			color.Printf(subjectMinorNameFormat, 1, " ", general.SuccessText("-"), general.FgBlue(subjectMinorName))
			color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 2, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.DockerMirrorsConfigFile))
			if err := general.WriteFile(cli.DockerMirrorsConfigFile, cli.DockerMirrorsConfig); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				fmt.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
		}
		// 配置 frpc
		if frpcFlag {
			subjectName = "frpc"
			descriptorText = "restart timing"
			color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlue(subjectName))
			color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.FrpcConfigFile))
			if err := general.WriteFile(cli.FrpcConfigFile, cli.FrpcConfig); err != nil {
				color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				fmt.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
		}
		// 配置 git
		if gitFlag {
			subjectName = "git"
			descriptorText = "configuration file"
			color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlue(subjectName))
			color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.GitConfigFile))
			if err := general.WriteFile(cli.GitConfigFile, cli.GitConfig); err != nil {
				color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				fmt.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
		}
		// 配置 golang
		if goFlag {
			subjectName = "go"
			descriptorText = "environment file"
			color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlue(subjectName))
			color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.GoConfigFile))
			if err := general.WriteFile(cli.GoConfigFile, cli.GoConfig); err != nil {
				color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				fmt.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
		}
		// 配置 pip
		if pipFlag {
			subjectName = "pip"
			descriptorText = "mirrors"
			color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlue(subjectName))
			color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.PipConfigFile))
			if err := general.WriteFile(cli.PipConfigFile, cli.PipConfig); err != nil {
				color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				fmt.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
		}
		// 配置 system-checkupdates
		if systemcheckupdatesFlag {
			subjectName = "system-checkupdates"
			color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlue(subjectName))
			// system-checkupdates timer
			subjectMinorName = "system-checkupdates timer"
			descriptorText = "timer"
			color.Printf(subjectMinorNameFormat, 1, " ", general.SuccessText("-"), general.FgBlue(subjectMinorName))
			color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.SystemCheckupdatesTimerConfigFile))
			if err := general.WriteFile(cli.SystemCheckupdatesTimerConfigFile, cli.SystemCheckupdatesTimerConfig); err != nil {
				errorFormat = "%*s%s %s: %s\n"
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				successFormat = "%*s%s %s: %s\n"
				fmt.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
			}
			// system-checkupdates service
			subjectMinorName = "system-checkupdates service"
			descriptorText = "service"
			color.Printf(subjectMinorNameFormat, 1, " ", general.SuccessText("-"), general.FgBlue(subjectMinorName))
			color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
			color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(cli.SystemCheckupdatesServiceConfigFile))
			if err := general.WriteFile(cli.SystemCheckupdatesServiceConfigFile, cli.SystemCheckupdatesServiceConfig); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				fmt.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellow("Setup completed"))
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
