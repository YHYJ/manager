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
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Operate configuration file",
	Long:  `Manipulate the program's configuration files, including generating and printing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取配置文件路径
		cfgFile, _ := cmd.Flags().GetString("config")
		// 解析参数
		createFlag, _ := cmd.Flags().GetBool("create")
		forceFlag, _ := cmd.Flags().GetBool("force")
		printFlag, _ := cmd.Flags().GetBool("print")

		// 创建配置文件流程
		if createFlag {
			cli.CreateConfigFile(cfgFile, forceFlag)
		}

		// 打印配置文件流程
		if printFlag {
			cli.PrintConfigFile(cfgFile)
		}
	},
}

func init() {
	configCmd.Flags().BoolP("create", "", false, "Create a default configuration file")
	configCmd.Flags().BoolP("force", "", false, "Overwrite existing configuration files")
	configCmd.Flags().BoolP("print", "", false, "Print configuration file content")

	configCmd.Flags().BoolP("help", "h", false, "help for config command")
	rootCmd.AddCommand(configCmd)
}
