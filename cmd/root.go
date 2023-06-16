/*
File: root.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-07 16:06:45

Description: 程序未带子命令或参数时执行
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "manager",
	Short: "Self-developed program manager",
	Long:  `Manager for self-developed programs, including installation and configuration.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile = "/etc/manager/config.toml"

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "help for Manager")

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", cfgFile, "Config file")
}
