/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 13:35:06

Description: 程序子命令'install'时执行
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install or update programs/scripts",
	Long:  `Install or update self-developed programs/scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("install called")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
