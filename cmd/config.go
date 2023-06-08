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
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Operate configuration file",
	Long:  `Manipulate the program's configuration files, including generating and printing.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
