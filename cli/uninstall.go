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
	"strings"

	"github.com/gookit/color"
	"github.com/yhyj/manager/general"
)

// UninstallSelf 卸载管理程序本身
//
// 参数：
//   - config: 解析 toml 配置文件得到的配置项
func UninstallSelf(config *general.Config) {
	// 程序文件
	program := config.Program.Self.Name // 程序名

	// 开始卸载提示
	color.Info.Tips("Uninstall \x1b[3m%s\x1b[0m programs", general.FgCyanText(program))

	// 检测主文件是否存在来决定是否在选项中显示
	programMainFile := filepath.Join(config.Program.ProgramPath, program) // 程序主文件路径
	if general.FileExist(programMainFile) {
		color.Printf("%s\n", strings.Repeat(general.Separator2st, general.SeparatorBaseLength))
	} else {
		color.Printf("%s\n", strings.Repeat(general.Separator3st, general.SeparatorBaseLength))
		color.Warn.Tips("Program \x1b[3m%s\x1b[0m is not installed", general.FgCyanText(program))
		return
	}

	// 设置文本参数
	textLength := 0 // 用于计算最后一行文本的长度，以便输出适当长度的分隔符

	// 确认是否要卸载
	question := color.Sprintf(general.UninstallTips, program)
	answer, err := general.AreYouSure(general.QuestionText(question), false)
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
		return
	}
	switch answer {
	case true:
		// 记账文件
		pocketDir := filepath.Join(config.Program.PocketPath, program)    // 记账文件夹路径
		pocketFile := filepath.Join(pocketDir, config.Program.PocketFile) // 记账文件路径
		pocketLines := make([]string, 0)                                  // 记账文件内容
		if general.FileExist(pocketFile) {                                // 读取记账文件内容
			pocketLines, err = general.ReadFile(pocketFile)
			if err != nil {
				fileName, lineNo := general.GetCallerInfo()
				color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
				return
			}
		}
		color.Printf("%s\n", strings.Repeat(general.Separator2st, len(question)))

		// 卸载程序
		for _, pocketLine := range pocketLines {
			if err := general.Uninstall(pocketLine); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				text := color.Sprintf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
				color.Print(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(general.DelayTime)      // 添加一个延时，使输出更加顺畅
				return
			}
		}

		// 删除记账文件
		if err := general.DeleteFile(pocketDir); err != nil {
			fileName, lineNo := general.GetCallerInfo()
			text := color.Sprintf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			color.Print(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(general.DelayTime)      // 添加一个延时，使输出更加顺畅
			return
		}

		// 本次卸载结束分隔符
		text := color.Sprintf("%s %s %s\n", general.SuccessFlag, general.FgGreenText(program), general.FgMagentaText("uninstalled"))
		color.Print(text)
		textLength = general.RealLength(text) // 分隔符长度

		// 分隔符和延时（延时使输出更加顺畅）
		general.PrintDelimiter(textLength) // 分隔符
		general.Delay(general.DelayTime)   // 添加一个延时，使输出更加顺畅
	case false:
		return
	default:
		color.Printf("%s\n", strings.Repeat(general.Separator3st, len(question)))
		color.Warn.Tips("%s: %s", "Unexpected answer", answer)
		return
	}
}

// Uninstall 卸载指定程序
//
// 参数：
//   - config: 解析 toml 配置文件得到的配置项
//   - category: 要卸载的类别，支持 uninstall 子命令除 '--all' 和 '--self' 之外的所有 Flags
func Uninstall(config *general.Config, category string) {
	// 从配置读取指定类别的程序名
	var programNames []string // 程序名切片
	switch category {
	case "go":
		programNames = config.Program.Go.Names
	case "shell":
		programNames = config.Program.Shell.Names
	default:
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s Category '%s' mismatch\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), category)
		return
	}

	// 检测程序主文件是否存在来决定是否在选项中显示
	installedPrograms := make([]string, 0) // 已安装程序
	for _, program := range programNames {
		programMainFile := filepath.Join(config.Program.ProgramPath, program) // 程序主文件路径
		if general.FileExist(programMainFile) {
			installedPrograms = append(installedPrograms, program)
		}
	}

	// 显示项排序
	sort.Strings(installedPrograms)

	// 开始卸载提示
	totalNum := len(programNames)          // 总程序数
	installedNum := len(installedPrograms) // 已安装程序数
	negatives := strings.Builder{}
	negatives.WriteString(color.Sprintf("%s Uninstall %s programs, %d/%d installed\n", general.InfoText("INFO:"), general.FgCyanText(category, "-based"), installedNum, totalNum))

	// 让用户选择需要卸载的程序
	selectedPrograms, err := general.MultipleSelectionFilter(installedPrograms, installedPrograms, negatives.String())
	if err != nil {
		fileName, lineNo := general.GetCallerInfo()
		color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
	}

	// 选择项排序
	sort.Strings(selectedPrograms)

	// 留屏信息
	if len(selectedPrograms) > 0 {
		negatives.WriteString(color.Sprintf("%s Selected: %s\n", general.InfoText("INFO:"), general.FgCyanText(strings.Join(selectedPrograms, ", "))))
		negatives.WriteString(color.Sprintf("%s", strings.Repeat(general.Separator1st, general.SeparatorBaseLength)))
		color.Println(negatives.String())
	}

	// 设置文本参数
	textLength := 0 // 用于计算最后一行文本的长度，以便输出适当长度的分隔符

	// 遍历所选程序/脚本名
	for _, program := range selectedPrograms {
		// 确认是否要卸载
		question := color.Sprintf(general.UninstallTips, program)
		answer, err := general.AreYouSure(general.QuestionText(question), false)
		if err != nil {
			fileName, lineNo := general.GetCallerInfo()
			color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
			return
		}
		switch answer {
		case true:
			color.Printf("%s\n", strings.Repeat(general.Separator2st, len(question)))
			// 记账文件
			pocketDir := filepath.Join(config.Program.PocketPath, program)    // 记账文件夹路径
			pocketFile := filepath.Join(pocketDir, config.Program.PocketFile) // 记账文件路径
			pocketLines := make([]string, 0)                                  // 记账文件内容
			if general.FileExist(pocketFile) {                                // 读取记账文件内容
				pocketLines, err = general.ReadFile(pocketFile)
				if err != nil {
					fileName, lineNo := general.GetCallerInfo()
					color.Printf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
					continue
				}
			}

			// 卸载程序
			for _, pocketLine := range pocketLines {
				if err := general.Uninstall(pocketLine); err != nil {
					fileName, lineNo := general.GetCallerInfo()
					text := color.Sprintf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
					color.Print(text)
					// 分隔符和延时（延时使输出更加顺畅）
					textLength = general.RealLength(text) // 分隔符长度
					general.PrintDelimiter(textLength)    // 分隔符
					general.Delay(general.DelayTime)      // 添加一个延时，使输出更加顺畅
					return
				}
			}

			// 删除记账文件
			if err := general.DeleteFile(pocketDir); err != nil {
				fileName, lineNo := general.GetCallerInfo()
				text := color.Sprintf("%s %s %s\n", general.DangerText(general.ErrorInfoFlag), general.SecondaryText("[", fileName, ":", lineNo+1, "]"), err)
				color.Print(text)
				// 分隔符和延时（延时使输出更加顺畅）
				textLength = general.RealLength(text) // 分隔符长度
				general.PrintDelimiter(textLength)    // 分隔符
				general.Delay(general.DelayTime)      // 添加一个延时，使输出更加顺畅
				return
			}

			// 本次卸载结束分隔符
			text := color.Sprintf("%s %s %s\n", general.SuccessFlag, general.FgGreenText(program), general.FgMagentaText("uninstalled"))
			color.Print(text)
			textLength = general.RealLength(text) // 分隔符长度

			// 分隔符和延时（延时使输出更加顺畅）
			general.PrintDelimiter(textLength) // 分隔符
			general.Delay(general.DelayTime)   // 添加一个延时，使输出更加顺畅
		case false:
			continue
		default:
			color.Printf("%s\n", strings.Repeat(general.Separator3st, len(general.UninstallTips)))
			color.Warn.Tips("%s: %s", "Unexpected answer", answer)
			continue
		}
	}
}
