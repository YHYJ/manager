/*
File: uninstall.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-27 12:35:06

Description: 执行子命令 'uninstall'
*/

package cmd

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall programs and scripts",
	Long:  `Uninstall self-developed programs and scripts.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Println("uninstall called")
	},
}

func init() {
	uninstallCmd.Flags().BoolP("self", "", false, "Uninstall itself (Can only be called alone)")
	uninstallCmd.Flags().BoolP("all", "", false, "Uninstall all programs and scripts")
	uninstallCmd.Flags().BoolP("go", "", false, "Uninstall programs developed based on go")
	uninstallCmd.Flags().BoolP("shell", "", false, "Uninstall shell scripts")

	uninstallCmd.Flags().BoolP("help", "h", false, "help for uninstall command")
	rootCmd.AddCommand(uninstallCmd)
}
