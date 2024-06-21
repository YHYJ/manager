/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:05:47

Description: 子命令 'config' 的实现
*/

package cli

import (
	"strings"

	"github.com/gookit/color"
	"github.com/yhyj/manager/general"
)

// CreateConfigFile 创建配置文件
//
// 参数：
//   - configFile: 配置文件路径
func CreateConfigFile(configFile string) {
	// 检查配置文件是否存在
	fileExist := general.FileExist(configFile)

	// 检测并创建配置文件
	if fileExist {
		// 询问是否覆写已存在的配置文件
		question := color.Sprintf(general.OverWriteTips, "Configuration")
		overWrite, err := general.AskUser(general.QuestionText(question), []string{"y", "N"})
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s -> Unable to get answers: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			return
		}

		switch overWrite {
		case "y":
			// 与用户交互获取配置信息
			general.InstallMethod, _ = general.AskUser(general.QuestionText(color.Sprintf(general.SelectOneTips, "the installation method")), general.AllInstallMethod)
			general.HttpProxy, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "HTTP_PROXY")), general.HttpProxy)
			general.HttpsProxy, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "HTTPS_PROXY")), general.HttpProxy)

			if err := general.DeleteFile(configFile); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Printf("%s %s -> Unable to delete file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
				return
			}
			if err := general.CreateFile(configFile); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Printf("%s %s -> Unable to create file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
				return
			}
			_, err := general.WriteTomlConfig(configFile)
			if err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Printf("%s %s -> Unable to write file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
				return
			}
			color.Printf("Create %s: %s\n", general.PrimaryText(configFile), general.SuccessText("file overwritten"))
		case "n":
			return
		default:
			color.Printf("%s\n", strings.Repeat(general.Separator3st, len(question)))
			color.Warn.Tips("%s: %s", "Unexpected answer", overWrite)
			return
		}
	} else {
		// 与用户交互获取代理配置
		general.HttpProxy, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "HTTP_PROXY")), general.HttpProxy)
		general.HttpsProxy, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "HTTPS_PROXY")), general.HttpProxy)

		if err := general.CreateFile(configFile); err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s -> Unable to create file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			return
		}
		_, err := general.WriteTomlConfig(configFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s -> Unable to write file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			return
		}
		color.Printf("Create %s: %s\n", general.PrimaryText(configFile), general.SuccessText("file created"))
	}
}

// OpenConfigFile 打开配置文件
//
// 参数：
//   - configFile: 配置文件路径
func OpenConfigFile(configFile string) {
	// 检查配置文件是否存在
	fileExist := general.FileExist(configFile)

	if fileExist {
		editor := general.GetVariable("EDITOR")
		if editor == "" {
			editor = "vim"
			err := general.RunCommand(editor, []string{configFile})
			if err != nil {
				editor = "vi"
				if err = general.RunCommand(editor, []string{configFile}); err != nil {
					fileName, lineNo := general.GetCallerInfo()
					color.Printf("%s %s -> Unable to open file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
				}
			}
		} else {
			if err := general.RunCommand(editor, []string{configFile}); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Printf("%s %s -> Unable to open file: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			}
		}
	}
}

// PrintConfigFile 打印配置文件内容
//
// 参数：
//   - configFile: 配置文件路径
func PrintConfigFile(configFile string) {
	// 检查配置文件是否存在
	fileExist := general.FileExist(configFile)

	var (
		configFileNotFoundMessage = "Configuration file not found (use --create to create a configuration file)" // 配置文件不存在
	)

	if fileExist {
		configTree, err := general.GetTomlConfig(configFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s -> Unable to get config: %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		} else {
			color.Println(general.PrimaryText(configTree))
		}
	} else {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s -> %s\n", general.DangerText("Error:"), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), configFileNotFoundMessage)
	}
}
