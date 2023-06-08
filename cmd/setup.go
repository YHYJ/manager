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
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure installed programs/scripts",
	Long:  `Configure installed self-developed programs/scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setup called")
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
