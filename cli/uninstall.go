/*
File: uninstall.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-27 12:32:16

Description: 子命令 'uninstall' 的实现
*/

package cli

import (
	"os"
	"path/filepath"

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
	name := config.Program.Self.Name                                // 程序名
	localProgram := filepath.Join(config.Program.ProgramPath, name) // 程序路径
	// 卸载程序文件
	if err := os.Remove(localProgram); err != nil {
		text := color.Sprintf("%s\n", general.ErrorText(err))
		color.Printf(text)
		// 分隔符和延时（延时使输出更加顺畅）
		textLength = general.RealLength(text) // 分隔符长度
		general.PrintDelimiter(textLength)    // 分隔符
		general.Delay(0.1)                    // 0.1s
		return
	}

	// 资源文件
	localResourcesDesktopFile := filepath.Join(config.Program.ResourcesPath, "applications", color.Sprintf("%s.desktop", name)) // desktop 文件
	localResourcesIconFolder := filepath.Join(config.Program.ResourcesPath, "pixmaps")                                          // icon 文件夹
	// 卸载资源文件 - desktop 文件
	if general.FileExist(localResourcesDesktopFile) {
		if err := general.Uninstall(localResourcesDesktopFile); err != nil {
			text := color.Sprintf("%s\n", general.ErrorText(err))
			color.Printf(text)
			// 分隔符和延时（延时使输出更加顺畅）
			textLength = general.RealLength(text) // 分隔符长度
			general.PrintDelimiter(textLength)    // 分隔符
			general.Delay(0.1)                    // 0.1s
			return
		}
	}
	// 卸载资源文件 - icon 文件
	if general.FileExist(localResourcesIconFolder) {
		color.Notice.Tips("TODO: 在包安装记录文件功能之后实现")
	}

	// 卸载自动补全脚本
	for _, completionDir := range config.Program.Go.CompletionDir {
		if general.FileExist(completionDir) {
			completionFile := filepath.Join(completionDir, name)
			if err := general.Uninstall(completionFile); err != nil {
				text := color.Sprintf("%s %s\n", general.ErrorFlag, general.ErrorText(general.AcsUninstallFailedMessage))
				color.Printf(text)
				textLength = general.RealLength(text) // 分隔符长度
				continue
			} else {
				text := color.Sprintf("%s %s\n", general.SuccessFlag, general.SecondaryText(general.AcsUninstallSuccessMessage))
				color.Printf(text)
				textLength = general.RealLength(text) // 分隔符长度
				break
			}
		}
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
}

// UninstallShellBasedProgram 卸载基于 Shell 的程序
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
func UninstallShellBasedProgram(configTree *toml.Tree) {
}
