//go:build windows

/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 14:13:47

Description: 子命令 'setup' 的实现
*/

package cli

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/yhyj/manager/general"
)

// ProgramConfigurator 程序配置器
// 参数：
//   - flags: 系统信息各部分的开关
func ProgramConfigurator(flags map[string]bool) {
	// 预定义变量
	var (
		home        = general.UserInfo.HomeDir
		userName, _ = general.GetUserName()
		name        = general.GetHostname()
		sep         = strings.Repeat(" ", 4)
	)

	var (
		subjectName    string
		descriptorText string
		writeMode      string = "t"
	)

	// 预定义输出格式
	var (
		descriptorFormat   = "%*s%s %s: %s %s %s\n"
		targetFileFormat   = "%*s%s %s: %s\n"
		askItemTitleFormat = "%*s%s %s:\n"
		askItemsFormat     = "%*s%s "
		errorFormat        = "%*s%s %s: %s\n\n"
		successFormat      = "%*s%s %s: %s\n\n"
	)

	// 配置 cobra
	if flags["cobraFlag"] {
		// 提示
		subjectName = "cobra-cli"
		descriptorText = "configuration file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// cobra 的依赖
			CobraDependencies = "cobra-cli"                        // 主程序
			CobraConfigFile   = filepath.Join(home, ".cobra.yaml") // 配置文件
			// cobra 配置
			cobraConfigFormat = "author: %s <%s>\nlicense: %s\nuseViper: %v\n"
			cobraAuthorName   = userName
			cobraAuthorEmail  = "email@example.com"
			cobraLicense      = "GPLv3"
			cobraUseViper     = "false"
		)

		// 检测
		if _, err := exec.LookPath(CobraDependencies); err != nil {
			color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.LightText("Target file"), general.CommentText(CobraConfigFile))
			// 创建配置文件
			if err := general.CreateFile(CobraConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Configuration"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				cobraAuthorName, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "authorName")), cobraAuthorName)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				cobraAuthorEmail, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "authorEmail")), cobraAuthorEmail)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				cobraLicense, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "license")), cobraLicense)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				cobraUseViper, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "useViper")), cobraUseViper)

				// 配置
				CobraConfigContent := color.Sprintf(cobraConfigFormat, cobraAuthorName, cobraAuthorEmail, cobraLicense, cobraUseViper)
				if err := general.WriteFile(CobraConfigFile, CobraConfigContent, writeMode); err != nil {
					color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
				} else {
					color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.SuccessFlag)
				}
			}
		}
	}

	// 配置 git
	if flags["gitFlag"] {
		// 提示
		subjectName = "git"
		descriptorText = "configuration file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// git 的依赖
			GitDependencies = "git"                             // 主程序
			GitConfigFile   = filepath.Join(home, ".gitconfig") // 配置文件
			// git 配置
			gitConfigFormat      = "[user]\n%sname = %s\n%semail = %s\n[core]\n%seditor = %s\n%sautocrlf = %s\n[merge]\n%stool = %s\n[color]\n%sui = %v\n[pull]\n%srebase = %v\n[filter \"lfs\"]\n%sclean = %s\n%ssmudge = %s\n%sprocess = %s\n%srequired = %v\n"
			gitUserEmail         = "email@example.com"
			gitCoreEditor        = "vim"
			gitCoreAutoCRLF      = "input"
			gitMergeTool         = "vimdiff"
			gitColorUI           = "true"
			gitPullRebase        = "false"
			gitFilterLfsClean    = "git-lfs clean -- %f"
			gitFilterLfsSmudge   = "git-lfs smudge -- %f"
			gitFilterLfsProcess  = "git-lfs filter-process"
			gitFilterLfsRequired = "true"
		)

		// 检测
		if _, err := exec.LookPath(GitDependencies); err != nil {
			color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.LightText("Target file"), general.CommentText(GitConfigFile))
			// 创建配置文件
			if err := general.CreateFile(GitConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Configuration"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				name, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[user].name")), name)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitUserEmail, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[user].email")), gitUserEmail)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitCoreEditor, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[core].editor")), gitCoreEditor)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitCoreAutoCRLF, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[core].autocrlf")), gitCoreAutoCRLF)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitMergeTool, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[merge].tool")), gitMergeTool)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitColorUI, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[color].ui")), gitColorUI)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitPullRebase, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[pull].rebase")), gitPullRebase)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitFilterLfsClean, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[filter \"lfs\"].clean")), gitFilterLfsClean)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitFilterLfsSmudge, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[filter \"lfs\"].smudge")), gitFilterLfsSmudge)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitFilterLfsProcess, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[filter \"lfs\"].process")), gitFilterLfsProcess)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitFilterLfsRequired, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[filter \"lfs\"].required")), gitFilterLfsRequired)

				// 配置
				GitConfigContent := color.Sprintf(gitConfigFormat, sep, name, sep, gitUserEmail, sep, gitCoreEditor, sep, gitCoreAutoCRLF, sep, gitMergeTool, sep, gitColorUI, sep, gitPullRebase, sep, gitFilterLfsClean, sep, gitFilterLfsSmudge, sep, gitFilterLfsProcess, sep, gitFilterLfsRequired)
				if err := general.WriteFile(GitConfigFile, GitConfigContent, writeMode); err != nil {
					color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
				} else {
					color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.SuccessFlag)
				}
			}
		}
	}

	// 配置 pip
	if flags["pipFlag"] {
		// 提示
		subjectName = "pip"
		descriptorText = "mirrors"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// pip 的依赖
			PipDependencies = "pip"                                 // 主程序
			PipConfigFile   = filepath.Join(home, "pip", "pip.ini") // 配置文件
			// pip 配置
			pipConfigFormat = "[global]\nindex-url = %s\ntrusted-host = %s\n"
			pipIndexUrl     = "https://mirrors.aliyun.com/pypi/simple"
		)

		// 检测
		if _, err := exec.LookPath(PipDependencies); err != nil {
			color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.LightText("Target file"), general.CommentText(PipConfigFile))
			// 创建配置文件
			if err := general.CreateFile(PipConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Configuration"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				pipIndexUrl, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "index-url")), pipIndexUrl)

				// 需要获取交互结果的配置项
				pipTrustedHost, _ := general.GetUrlHost(pipIndexUrl)

				// 配置
				PipConfigContent := color.Sprintf(pipConfigFormat, pipIndexUrl, pipTrustedHost)
				if err := general.WriteFile(PipConfigFile, PipConfigContent, writeMode); err != nil {
					color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
				} else {
					color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.SuccessFlag)
				}
			}
		}
	}
}
