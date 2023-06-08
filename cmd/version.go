/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 13:46:23

Description: 程序子命令'version'时执行
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print program version",
	Long:  `Print program version and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version called")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
