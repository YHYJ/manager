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
	"github.com/yhyj/manager/cli"
	"github.com/yhyj/manager/general"
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
		cfgFileExist := general.FileExist(cfgFile)

		// 执行配置文件操作
		if createFlag {
			if cfgFileExist {
				if forceFlag {
					general.DeleteFile(cfgFile)
					general.CreateFile(cfgFile)
					_, err := cli.WriteTomlConfig(cfgFile)
					if err != nil {
						fmt.Printf(general.ErrorBaseFormat, err)
						return
					}
					fmt.Printf(general.InfoPrefixSuffixFormat, "Create", " ", cfgFile, ": ", "file overwritten")
				} else {
					fmt.Printf(general.InfoPrefixSuffixFormat, "Create", " ", cfgFile, ": ", "file exists (use --force to overwrite)")
				}
			} else {
				if err := general.CreateFile(cfgFile); err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					return
				}
				_, err := cli.WriteTomlConfig(cfgFile)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
					return
				}
				fmt.Printf(general.InfoPrefixSuffixFormat, "Create", " ", cfgFile, ": ", "file created")
			}
		}

		if printFlag {
			if cfgFileExist {
				configTree, err := cli.GetTomlConfig(cfgFile)
				if err != nil {
					fmt.Printf(general.ErrorBaseFormat, err)
				} else {
					fmt.Println(configTree)
				}
			} else {
				fmt.Printf(general.ErrorBaseFormat, cfgFileNotFoundMessage)
			}
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
