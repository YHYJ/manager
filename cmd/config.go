/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-07 16:41:43

Description: 执行子命令 'config'
*/

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yhyj/manager/cli"
	"github.com/yhyj/manager/general"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Operate configuration file",
	Long:  `Manipulate the program's configuration files, including generating and printing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取配置文件路径
		configFile, _ := cmd.Flags().GetString("config")
		// 解析参数
		createFlag, _ := cmd.Flags().GetBool("create")
		openFlag, _ := cmd.Flags().GetBool("open")
		printFlag, _ := cmd.Flags().GetBool("print")

		// 检查参数
		if !createFlag && !printFlag && !openFlag {
			cmd.Help()
			general.Notifier = append(general.Notifier, "Please refer to the above help information")
		}

		// 创建配置文件流程
		if createFlag {
			cli.CreateConfigFile(configFile)
		}

		// 打开配置文件流程
		if openFlag {
			cli.OpenConfigFile(configFile)
		}

		// 打印配置文件流程
		if printFlag {
			cli.PrintConfigFile(configFile)
		}

		// 显示通知
		general.Notification()
	},
}

func init() {
	configCmd.Flags().Bool("create", false, "Create a default configuration file")
	configCmd.Flags().Bool("open", false, "Open the configuration file with the default editor")
	configCmd.Flags().Bool("print", false, "Print configuration file content")

	configCmd.Flags().BoolP("help", "h", false, "help for config command")
	rootCmd.AddCommand(configCmd)
}
