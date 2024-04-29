/*
File: uninstall.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-27 12:32:16

Description: 子命令 'uninstall' 的实现
*/

package cli

import (
	"path/filepath"
	"sort"

	"github.com/gookit/color"
	"github.com/pelletier/go-toml"
	"github.com/yhyj/manager/general"
)

// UninstallSelfProgram 卸载管理程序本身
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
func UninstallSelfProgram(configTree *toml.Tree) {
	// 获取配置项
	config, err := general.LoadConfigToStruct(configTree)
	if err != nil {
		color.Error.Println(err)
		return
	}

	// 开始卸载提示
	color.Info.Tips("Uninstall \x1b[3m%s\x1b[0m programs", general.FgCyanText(config.Program.Self.Name))

	// 设置文本参数
	textLength := 0 // 用于计算最后一行文本的长度，以便输出适当长度的分隔符

	// 程序文件
	name := config.Program.Self.Name // 程序名

	// 记账文件
	pocketDir := filepath.Join(config.Program.PocketPath, name)       // 记账文件夹路径
	pocketFile := filepath.Join(pocketDir, config.Program.PocketFile) // 记账文件路径
	pocketLines, err := general.ReadFile(pocketFile)                  // 读取记账文件内容
	if err != nil {
		color.Error.Println(err)
		return
	}

	// 确认是否要卸载
	question := general.QuestionText("Do you plan to delete these software?")
	answer, err := general.AskUser(question, "y/N")
	if err != nil {
		color.Error.Println(err)
		return
	}
	if answer == "n" {
		return
	}

	// 卸载程序
	for _, pocketLine := range pocketLines {
		if err := general.Uninstall(pocketLine); err != nil {
			text := color.Sprintf("%s\n", general.ErrorText(err))
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}
	}

	// 删除记账文件
	if err := general.DeleteFile(pocketDir); err != nil {
		text := color.Sprintf("%s\n", general.ErrorText(err))
		color.Printf(text)
		// 分隔符和延时（延时使输出更加顺畅）
		textLength = general.RealLength(text) // 分隔符长度
		general.PrintDelimiter(textLength)    // 分隔符
		general.Delay(0.1)                    // 0.1s
		return
	}

	// 本次卸载结束分隔符
	text := color.Sprintf("%s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgMagentaText("uninstalled"))
	color.Printf(text)
	textLength = general.RealLength(text) // 分隔符长度

	// 分隔符和延时（延时使输出更加顺畅）
	general.PrintDelimiter(textLength) // 分隔符
	general.Delay(0.1)                 // 0.01s
}

// UninstallGolangBasedProgram 卸载基于 Golang 的程序
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
func UninstallGolangBasedProgram(configTree *toml.Tree) {
	// 获取配置项
	config, err := general.LoadConfigToStruct(configTree)
	if err != nil {
		color.Error.Println(err)
		return
	}

	// 开始卸载提示
	color.Info.Tips("Uninstall \x1b[3m%s\x1b[0m programs", general.FgCyanText("golang-based"))

	// 设置文本参数
	textLength := 0 // 用于计算最后一行文本的长度，以便输出适当长度的分隔符

	// 让用户选择需要卸载的程序
	selectedNames, err := general.MultipleSelectionFilter(config.Program.Go.Names)
	if err != nil {
		color.Error.Println(err)
	}
	// 对所选的程序进行排序
	sort.Strings(selectedNames)

	// 确认是否要卸载
	question := general.QuestionText("Do you plan to delete these software?")
	answer, err := general.AskUser(question, "y/N")
	if err != nil {
		color.Error.Println(err)
		return
	}
	if answer == "n" {
		return
	}

	// 遍历所选脚本名
	for _, name := range selectedNames {
		// 记账文件
		pocketDir := filepath.Join(config.Program.PocketPath, name)       // 记账文件夹路径
		pocketFile := filepath.Join(pocketDir, config.Program.PocketFile) // 记账文件路径
		pocketLines, err := general.ReadFile(pocketFile)                  // 读取记账文件内容
		if err != nil {
			color.Error.Println(err)
			return
		}

		// 卸载程序
		for _, pocketLine := range pocketLines {
			if err := general.Uninstall(pocketLine); err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
		}

		// 删除记账文件
		if err := general.DeleteFile(pocketDir); err != nil {
			text := color.Sprintf("%s\n", general.ErrorText(err))
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}

		// 本次卸载结束分隔符
		text := color.Sprintf("%s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgMagentaText("uninstalled"))
		color.Printf(text)
		textLength = general.RealLength(text) // 分隔符长度

		// 分隔符和延时（延时使输出更加顺畅）
		general.PrintDelimiter(textLength) // 分隔符
		general.Delay(0.1)                 // 0.01s
	}
}

// UninstallShellBasedProgram 卸载基于 Shell 的程序
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
func UninstallShellBasedProgram(configTree *toml.Tree) {
	// 获取配置项
	config, err := general.LoadConfigToStruct(configTree)
	if err != nil {
		color.Error.Println(err)
		return
	}

	// 开始卸载提示
	color.Info.Tips("Uninstall \x1b[3m%s\x1b[0m programs", general.FgCyanText("shell-based"))

	// 设置文本参数
	textLength := 0 // 用于计算最后一行文本的长度，以便输出适当长度的分隔符

	// 让用户选择需要卸载的程序
	selectedNames, err := general.MultipleSelectionFilter(config.Program.Shell.Names)
	if err != nil {
		color.Error.Println(err)
	}
	// 对所选的程序进行排序
	sort.Strings(selectedNames)

	// 确认是否要卸载
	question := general.QuestionText("Do you plan to delete these software?")
	answer, err := general.AskUser(question, "y/N")
	if err != nil {
		color.Error.Println(err)
		return
	}
	if answer == "n" {
		return
	}

	// 遍历所选脚本名
	for _, name := range selectedNames {
		// 记账文件
		pocketDir := filepath.Join(config.Program.PocketPath, name)       // 记账文件夹路径
		pocketFile := filepath.Join(pocketDir, config.Program.PocketFile) // 记账文件路径
		pocketLines, err := general.ReadFile(pocketFile)                  // 读取记账文件内容
		if err != nil {
			color.Error.Println(err)
			return
		}

		// 卸载程序
		for _, pocketLine := range pocketLines {
			if err := general.Uninstall(pocketLine); err != nil {
				text := color.Sprintf("%s\n", general.ErrorText(err))
				color.Printf(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(0.1)                    // 0.1s
				return
			}
		}

		// 删除记账文件
		if err := general.DeleteFile(pocketDir); err != nil {
			text := color.Sprintf("%s\n", general.ErrorText(err))
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}

		// 本次卸载结束分隔符
		text := color.Sprintf("%s %s %s\n", general.SuccessFlag, general.FgGreenText(name), general.FgMagentaText("uninstalled"))
		color.Printf(text)
		textLength = general.RealLength(text) // 分隔符长度

		// 分隔符和延时（延时使输出更加顺畅）
		general.PrintDelimiter(textLength) // 分隔符
		general.Delay(0.1)                 // 0.01s
	}
}
