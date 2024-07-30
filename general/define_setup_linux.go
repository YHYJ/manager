//go:build linux

/*
File: define_setup_linux.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-07-30 14:53:39

Description: 供 'setup' 子命令使用的函数
*/

package general

import (
	"path/filepath"

	"github.com/gookit/color"
)

// ---------- 配置项

// Golang
var (
	// go 的依赖
	GolangDependencies = "go"                                        // 主程序
	GolangConfigFile   = filepath.Join(home, ".config", "go", "env") // 配置文件
	// go 配置
	golangConfigFormat = "GO111MODULE=%s\nGOBIN=%s\nGOPATH=%s\nGOCACHE=%s\nGOMODCACHE=%s\n"
	golangGO111MODULE  = "on"
	golangGOPATH       = filepath.Join(home, ".go")
	golangGOCACHE      = filepath.Join(home, ".cache", "go", "go-build")
	golangGOMODCACHE   = filepath.Join(home, ".cache", "go", "pkg", "mod")
)

// Pip
var (
	// pip 的依赖
	PipDependencies = "pip"                                             // 主程序
	PipConfigFile   = filepath.Join(home, ".config", "pip", "pip.conf") // 配置文件
	// pip 配置
	pipConfigFormat = "[global]\nindex-url = %s\ntrusted-host = %s\n"
	pipIndexUrl     = "https://mirrors.aliyun.com/pypi/simple"
)

// rebirth 服务配置完成后的重载/启用/重启
//
// 参数：
//   - name: 服务名称
//   - margin: 对齐时的边距
//   - coefficient: 边距应乘的系数
func rebirth(name string, margin, coefficient int) {
	cMargin := margin * coefficient

	// 重新加载 systemd 管理器配置
	reloadArgs := []string{"daemon-reload"} // 重载服务配置
	if _, _, err := RunCommandToBuffer("systemctl", reloadArgs); err != nil {
		color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), ErrorFlag, DangerText(err))
	}

	// 询问是否需要启用/重启服务
	checkStatusArgs := []string{"is-enabled", name} // 检测服务启用状态
	status, _, _ := RunCommandToBuffer("systemctl", checkStatusArgs)
	switch status {
	case "enabled":
		color.Printf(askItemsFormat, cMargin, " ", SuccessText("-"))

		question := color.Sprintf(RestartServiceTips, name)
		restart, err := AskUser(QuestionText(question), []string{"y", "N"})
		if err != nil {
			color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), ErrorFlag, DangerText("Unable to get answers: ", err))
		}
		switch restart {
		case "y":
			// 重启服务
			restartArgs := []string{"restart", name}
			if _, stderr, err := RunCommandToBuffer("systemctl", restartArgs); err != nil {
				color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), ErrorFlag, DangerText(stderr))
			} else {
				color.Printf(yesResultFormat, cMargin, " ", SuccessText("-"), SuccessFlag)
			}
		case "n":
			return
		default:
			color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), WarningFlag, WarnText("Unexpected answer: ", restart))
		}
	case "disabled":
		color.Printf(askItemsFormat, cMargin, " ", SuccessText("-"))

		question := color.Sprintf(EnableServiceTips, name)
		enable, err := AskUser(QuestionText(question), []string{"y", "N"})
		if err != nil {
			color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), ErrorFlag, DangerText("Unable to get answers: ", err))
		}
		switch enable {
		case "y":
			// 启用服务（并立即运行）
			enableArgs := []string{"enable", "--now", name}
			if _, stderr, err := RunCommandToBuffer("systemctl", enableArgs); err != nil {
				color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), ErrorFlag, DangerText(stderr))
			} else {
				color.Printf(yesResultFormat, cMargin, " ", SuccessText("-"), SuccessFlag)
			}
		case "n":
			return
		default:
			color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), WarningFlag, WarnText("Unexpected answer: ", enable))
		}
	case "static":
		color.Printf(yesResultFormat, cMargin, " ", SuccessText("-"), SuccessFlag)
	case "not-found":
		notFound := color.Sprintf(NotFoundServiceTips, name)
		color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), ErrorFlag, DangerText(notFound))
	}
}
