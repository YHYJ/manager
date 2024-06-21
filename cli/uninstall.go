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
	"strings"

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
		fileName, lineNo := general.GetCallerInfo()
		color.Danger.Printf("Load config error (%s:%d): %s\n", fileName, lineNo+1, err)
		return
	}

	// 程序文件
	name := config.Program.Self.Name // 程序名

	// 开始卸载提示
	color.Info.Tips("Uninstall \x1b[3m%s\x1b[0m programs", general.FgCyanText(name))

	// 检测主文件是否存在来决定是否在选项中显示
	programMainFile := filepath.Join(config.Program.ProgramPath, name) // 程序主文件路径
	if general.FileExist(programMainFile) {
		color.Printf("%s\n", strings.Repeat(general.Separator2st, general.SeparatorBaseLength))
	} else {
		color.Printf("%s\n", strings.Repeat(general.Separator3st, general.SeparatorBaseLength))
		color.Warn.Tips("Program \x1b[3m%s\x1b[0m is not installed", general.FgCyanText(name))
		return
	}

	// 记账文件
	pocketDir := filepath.Join(config.Program.PocketPath, name)       // 记账文件夹路径
	pocketFile := filepath.Join(pocketDir, config.Program.PocketFile) // 记账文件路径
	pocketLines := make([]string, 0)                                  // 记账文件内容
	if general.FileExist(pocketFile) {                                // 读取记账文件内容
		pocketLines, err = general.ReadFile(pocketFile)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Danger.Printf("Read file error (%s:%d): %s\n", fileName, lineNo+1, err)
			return
		}
	}

	// 确认是否要卸载
	answer, err := general.AskUser(general.QuestionText(general.UninstallTips), []string{"y", "N"})
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Danger.Printf("Ask user error (%s:%d): %s\n", fileName, lineNo+1, err)
		return
	}
	switch answer {
	case "y":
		color.Printf("%s\n", strings.Repeat(general.Separator2st, len(general.UninstallTips)))
	case "n":
		return
	default:
		color.Printf("%s\n", strings.Repeat(general.Separator3st, len(general.UninstallTips)))
		color.Warn.Tips("%s", "Unexpected answer")
		return
	}

	// 设置文本参数
	textLength := 0 // 用于计算最后一行文本的长度，以便输出适当长度的分隔符

	// 卸载程序
	for _, pocketLine := range pocketLines {
		if err := general.Uninstall(pocketLine); err != nil {
			fileName, lineNo := general.GetCallerInfo()
			text := color.Danger.Sprintf("Uninstall error (%s:%d): %s\n", fileName, lineNo+1, err)
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
		fileName, lineNo := general.GetCallerInfo()
		text := color.Danger.Sprintf("Uninstall error (%s:%d): %s\n", fileName, lineNo+1, err)
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
		fileName, lineNo := general.GetCallerInfo()
		color.Danger.Printf("Load config error (%s:%d): %s\n", fileName, lineNo+1, err)
		return
	}

	// 检测主文件是否存在来决定是否在选项中显示
	installedPrograms := make([]string, 0) // 已安装程序
	for _, program := range config.Program.Go.Names {
		programMainFile := filepath.Join(config.Program.ProgramPath, program) // 程序主文件路径
		if general.FileExist(programMainFile) {
			installedPrograms = append(installedPrograms, program)
		}
	}

	// 开始卸载提示
	totalNum := len(config.Program.Shell.Names) // 总程序数
	installedNum := len(installedPrograms)      // 已安装程序数
	negatives := strings.Builder{}
	negatives.WriteString(color.Sprintf("%s Uninstall \x1b[3m%s\x1b[0m programs, %d/%d installed\n", general.InfoText("INFO:"), general.FgCyanText("golang-based"), installedNum, totalNum))

	// 让用户选择需要卸载的程序
	selectedPrograms, err := general.MultipleSelectionFilter(installedPrograms, negatives.String())
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Danger.Printf("Filter error (%s:%d): %s\n", fileName, lineNo+1, err)
	}

	// 确认是否要卸载
	if len(selectedPrograms) != 0 {
		answer, err := general.AskUser(general.QuestionText(general.UninstallTips), []string{"y", "N"})
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Danger.Printf("Ask user error (%s:%d): %s\n", fileName, lineNo+1, err)
			return
		}
		switch answer {
		case "y":
			color.Printf("%s\n", strings.Repeat(general.Separator2st, len(general.UninstallTips)))
		case "n":
			return
		default:
			color.Printf("%s\n", strings.Repeat(general.Separator3st, len(general.UninstallTips)))
			color.Warn.Tips("%s", "Unexpected answer")
			return
		}
	}

	// 设置文本参数
	textLength := 0 // 用于计算最后一行文本的长度，以便输出适当长度的分隔符

	// 遍历所选程序名
	for _, program := range selectedPrograms {
		// 记账文件
		pocketDir := filepath.Join(config.Program.PocketPath, program)    // 记账文件夹路径
		pocketFile := filepath.Join(pocketDir, config.Program.PocketFile) // 记账文件路径
		pocketLines := make([]string, 0)                                  // 记账文件内容
		if general.FileExist(pocketFile) {                                // 读取记账文件内容
			pocketLines, err = general.ReadFile(pocketFile)
			if err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Danger.Printf("Read file error (%s:%d): %s\n", fileName, lineNo+1, err)
				continue
			}
		}

		// 卸载程序
		for _, pocketLine := range pocketLines {
			if err := general.Uninstall(pocketLine); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				text := color.Danger.Sprintf("Uninstall error (%s:%d): %s\n", fileName, lineNo+1, err)
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
			fileName, lineNo := general.GetCallerInfo()
			text := color.Danger.Sprintf("Uninstall error (%s:%d): %s\n", fileName, lineNo+1, err)
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}

		// 本次卸载结束分隔符
		text := color.Sprintf("%s %s %s\n", general.SuccessFlag, general.FgGreenText(program), general.FgMagentaText("uninstalled"))
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
		fileName, lineNo := general.GetCallerInfo()
		color.Danger.Printf("Load config error (%s:%d): %s\n", fileName, lineNo+1, err)
		return
	}

	// 检测主文件是否存在来决定是否在选项中显示
	installedPrograms := make([]string, 0) // 已安装程序
	for _, program := range config.Program.Shell.Names {
		programMainFile := filepath.Join(config.Program.ProgramPath, program) // 程序主文件路径
		if general.FileExist(programMainFile) {
			installedPrograms = append(installedPrograms, program)
		}
	}

	// 开始卸载提示
	totalNum := len(config.Program.Shell.Names) // 总程序数
	installedNum := len(installedPrograms)      // 已安装程序数
	negatives := strings.Builder{}
	negatives.WriteString(color.Sprintf("%s Uninstall \x1b[3m%s\x1b[0m programs, %d/%d installed\n", general.InfoText("INFO:"), general.FgCyanText("shell-based"), installedNum, totalNum))

	// 让用户选择需要卸载的程序
	selectedPrograms, err := general.MultipleSelectionFilter(installedPrograms, negatives.String())
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Danger.Printf("Filter error (%s:%d): %s\n", fileName, lineNo+1, err)
	}

	// 确认是否要卸载
	if len(selectedPrograms) != 0 {
		answer, err := general.AskUser(general.QuestionText(general.UninstallTips), []string{"y", "N"})
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Danger.Printf("Ask user error (%s:%d): %s\n", fileName, lineNo+1, err)
			return
		}
		switch answer {
		case "y":
			color.Printf("%s\n", strings.Repeat(general.Separator2st, len(general.UninstallTips)))
		case "n":
			return
		default:
			color.Printf("%s\n", strings.Repeat(general.Separator3st, len(general.UninstallTips)))
			color.Warn.Tips("%s", "Unexpected answer")
			return
		}
	}

	// 设置文本参数
	textLength := 0 // 用于计算最后一行文本的长度，以便输出适当长度的分隔符

	// 遍历所选脚本名
	for _, program := range selectedPrograms {
		// 记账文件
		pocketDir := filepath.Join(config.Program.PocketPath, program)    // 记账文件夹路径
		pocketFile := filepath.Join(pocketDir, config.Program.PocketFile) // 记账文件路径
		pocketLines := make([]string, 0)                                  // 记账文件内容
		if general.FileExist(pocketFile) {                                // 读取记账文件内容
			pocketLines, err = general.ReadFile(pocketFile)
			if err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Danger.Printf("Read file error (%s:%d): %s\n", fileName, lineNo+1, err)
				continue
			}
		}

		// 卸载程序
		for _, pocketLine := range pocketLines {
			if err := general.Uninstall(pocketLine); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				text := color.Danger.Sprintf("Uninstall error (%s:%d): %s\n", fileName, lineNo+1, err)
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
			fileName, lineNo := general.GetCallerInfo()
			text := color.Danger.Sprintf("Uninstall error (%s:%d): %s\n", fileName, lineNo+1, err)
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}

		// 本次卸载结束分隔符
		text := color.Sprintf("%s %s %s\n", general.SuccessFlag, general.FgGreenText(program), general.FgMagentaText("uninstalled"))
		color.Printf(text)
		textLength = general.RealLength(text) // 分隔符长度

		// 分隔符和延时（延时使输出更加顺畅）
		general.PrintDelimiter(textLength) // 分隔符
		general.Delay(0.1)                 // 0.01s
	}
}
