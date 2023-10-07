/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023]-06-07 16:41:43

Description: 程序子命令'config'时执行
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yhyj/manager/function"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Operate configuration file",
	Long:  `Manipulate the program's configuration files, including generating and printing.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 解析参数
		cfgFile, _ := cmd.Flags().GetString("config")
		createFlag, _ := cmd.Flags().GetBool("create")
		forceFlag, _ := cmd.Flags().GetBool("force")
		printFlag, _ := cmd.Flags().GetBool("print")

		var (
			cfgFileNotFoundMessage = "Configuration file not found (use --create to create a configuration file)" // 配置文件不存在
		)

		// 检查配置文件是否存在
		cfgFileExist := function.FileExist(cfgFile)

		// 执行配置文件操作
		if createFlag {
			if cfgFileExist {
				if forceFlag {
					function.DeleteFile(cfgFile)
					function.CreateFile(cfgFile)
					function.WriteTomlConfig(cfgFile)
					fmt.Printf("Create \x1b[33;1m%s\x1b[0m: file overwritten\n", cfgFile)
				} else {
					fmt.Printf("Create \x1b[33m%s\x1b[0m: file exists (use --force to overwrite)\n", cfgFile)
				}
			} else {
				err := function.CreateFile(cfgFile)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				_, err = function.WriteTomlConfig(cfgFile)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
					return
				}
				fmt.Printf("Create \x1b[33;1m%s\x1b[0m: file created\n", cfgFile)
			}
		}

		if printFlag {
			if cfgFileExist {
				configTree, err := function.GetTomlConfig(cfgFile)
				if err != nil {
					fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
				} else {
					fmt.Println(configTree)
				}
			} else {
				fmt.Printf("\x1b[31m%s\x1b[0m\n", cfgFileNotFoundMessage)
			}
		}
	},
}

func init() {
	configCmd.Flags().BoolP("create", "", false, "Create a default configuration file")
	configCmd.Flags().BoolP("force", "", false, "Overwrite existing configuration files")
	configCmd.Flags().BoolP("print", "", false, "Print configuration file content")

	configCmd.Flags().BoolP("help", "h", false, "help for config")
	rootCmd.AddCommand(configCmd)
}
