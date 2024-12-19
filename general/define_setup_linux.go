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
	"os/exec"
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
//   - owner: 服务所属用户（system 或 user）
//   - margin: 对齐时的边距
//   - coefficient: 边距应乘的系数
func rebirth(name string, owner string, margin, coefficient int) {
	manager := "--system"
	if owner == "user" {
		manager = "--user"
	}

	// 打印格式
	cMargin := margin * coefficient

	// 重新加载 systemd 管理器配置
	reloadArgs := []string{manager, "daemon-reload"} // 重载服务配置
	if _, _, err := RunCommandToBuffer("systemctl", reloadArgs); err != nil {
		color.Printf(noResultFormat, cMargin, " ", SuccessText("-"), ErrorFlag, DangerText(err))
	}

	// 询问是否需要启用/重启服务
	checkStatusArgs := []string{manager, "is-enabled", name} // 检测服务启用状态
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
			restartArgs := []string{manager, "restart", name}
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
			enableArgs := []string{manager, "enable", "--now", name}
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

// SetupUpdateChecker 配置 UpdateChecker
func SetupUpdateChecker() {
	// 提示
	subjectName = "update-checker"
	owner := "user"
	color.Printf("%s %s\n", SuccessText("==>"), FgBlueText(subjectName))

	// 检测
	if _, err := exec.LookPath(UpdateCheckerDependencies); err != nil {
		color.Printf(statusFormat, 2, " ", SuccessText("-"), LightText("Status"), NoticeText(color.Sprintf(InstallTips, subjectName)))
	} else {
		// ---------- Timer
		subjectMinorName = "timer"
		descriptorText = "timer"
		color.Printf(subjectMinorNameFormat, 2, " ", SuccessText("-"), FgBlueText(subjectMinorName))
		color.Printf(descriptorFormat, 4, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))
		color.Printf(targetFileFormat, 4, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(timerConfigFile))
		// 创建配置文件
		if err := CreateFile(timerConfigFile); err != nil {
			color.Printf(errorFormat, 4, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 4, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			timerDescription, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Unit].Description")), timerDescription)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			timerOnBootSec, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Timer].OnBootSec")), timerOnBootSec)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			timerOnUnitInactiveSec, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Timer].OnUnitInactiveSec")), timerOnUnitInactiveSec)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			timerAccuracySec, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Timer].AccuracySec")), timerAccuracySec)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			timerPersistent, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Timer].Persistent")), timerPersistent)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			timerWantedBy, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Install].WantedBy")), timerWantedBy)

			// 配置
			SystemCheckupdatesTimerConfigContent := color.Sprintf(timerConfigFormat, timerDescription, timerOnBootSec, timerOnUnitInactiveSec, timerAccuracySec, timerPersistent, timerWantedBy)
			if err := WriteFile(timerConfigFile, SystemCheckupdatesTimerConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 6, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 6, " ", SuccessText("-"), SuccessFlag)
			}

			// 重载服务
			color.Printf(askItemTitleFormat, 4, " ", SuccessText("-"), LightText("Rebirth"))
			name := color.Sprintf("%s.%s", subjectName, subjectMinorName)
			rebirth(name, owner, 2, 3)
		}

		// ---------- Service
		subjectMinorName = "service"
		descriptorText = "service"
		color.Printf(subjectMinorNameFormat, 2, " ", SuccessText("-"), FgBlueText(subjectMinorName))
		color.Printf(descriptorFormat, 4, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))
		color.Printf(targetFileFormat, 4, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(serviceConfigFile))
		// 创建配置文件
		if err := CreateFile(serviceConfigFile); err != nil {
			color.Printf(errorFormat, 4, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 4, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			serviceDescription, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Unit].Description")), serviceDescription)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			serviceAfter, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Unit].After")), serviceAfter)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			serviceWants, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Unit].Wants")), serviceWants)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			serviceType, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Service].Type")), serviceType)
			color.Printf(askItemsFormat, 6, " ", SuccessText("-"))
			serviceExecStart, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[Service].ExecStart")), serviceExecStart)

			// 配置
			SystemCheckupdatesServiceConfigContent := color.Sprintf(serviceConfigFormat, serviceDescription, serviceAfter, serviceWants, serviceType, serviceExecStart)
			if err := WriteFile(serviceConfigFile, SystemCheckupdatesServiceConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 6, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 6, " ", SuccessText("-"), SuccessFlag)
			}

			// 重载服务
			color.Printf(askItemTitleFormat, 4, " ", SuccessText("-"), LightText("Rebirth"))
			name := color.Sprintf("%s.%s", subjectName, subjectMinorName)
			rebirth(name, owner, 2, 3)
		}
	}
}

// SetupDocker 配置 Docker
func SetupDocker() {
	// 提示
	subjectName = "docker"
	owner := "system"
	descriptorText = "daemon configuration"
	color.Printf("%s %s\n", SuccessText("==>"), FgBlueText(subjectName))
	color.Printf(descriptorFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))

	// 检测
	if dockerdAbsPath, err := exec.LookPath(DockerDependencies); err != nil {
		color.Printf(statusFormat, 2, " ", SuccessText("-"), LightText("Status"), NoticeText(color.Sprintf(InstallTips, subjectName)))
	} else {
		color.Printf(targetFileFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(DockerServiceConfigFile))
		// 创建配置文件
		if err := CreateFile(DockerServiceConfigFile); err != nil {
			color.Printf(errorFormat, 2, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			dockerServiceDataRoot, _ = GetInput(QuestionText(color.Sprintf(InputTips, "--data-root")), dockerServiceDataRoot)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			HttpProxy, _ = GetInput(QuestionText(color.Sprintf(InputTips, "HTTP_PROXY")), HttpProxy)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			HttpsProxy, _ = GetInput(QuestionText(color.Sprintf(InputTips, "HTTPS_PROXY")), HttpProxy)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			noProxy, _ = GetInput(QuestionText(color.Sprintf(InputTips, "NO_PROXY")), noProxy)

			// 需要获取交互结果的配置项
			dockerHttpProxy := color.Sprintf("HTTP_PROXY=%s", HttpProxy)
			dockerHttpsProxy := color.Sprintf("HTTPS_PROXY=%s", HttpsProxy)
			dockerNoProxy := color.Sprintf("NO_PROXY=%s", noProxy)

			// 配置
			DockerServiceConfigContent := color.Sprintf(dockerServiceConfigFormat, dockerHttpProxy, dockerHttpsProxy, dockerNoProxy, dockerdAbsPath, dockerServiceDataRoot)
			if err := WriteFile(DockerServiceConfigFile, DockerServiceConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 4, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 4, " ", SuccessText("-"), SuccessFlag)
			}

			// 重载服务
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Rebirth"))
			rebirth(subjectName, owner, 2, 2)
		}
	}
}

// SetupFrpc 配置 Frpc
func SetupFrpc() {
	// 提示
	subjectName = "frpc"
	owner := "system"
	descriptorText = "restart timing"
	color.Printf("%s %s\n", SuccessText("==>"), FgBlueText(subjectName))
	color.Printf(descriptorFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))

	// 检测
	if _, err := exec.LookPath(FrpcDependencies); err != nil {
		color.Printf(statusFormat, 2, " ", SuccessText("-"), LightText("Status"), NoticeText(color.Sprintf(InstallTips, subjectName)))
	} else {
		color.Printf(targetFileFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(FrpcConfigFile))
		// 创建配置文件
		if err := CreateFile(FrpcConfigFile); err != nil {
			color.Printf(errorFormat, 2, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			frpcRestart, _ = GetInput(QuestionText(color.Sprintf(InputTips, "Restart")), frpcRestart)

			// 配置
			FrpcConfigContent := color.Sprintf(frpcConfigFormat, frpcRestart)
			if err := WriteFile(FrpcConfigFile, FrpcConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 4, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 4, " ", SuccessText("-"), SuccessFlag)
			}

			// 重载服务
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Rebirth"))
			rebirth(subjectName, owner, 2, 2)
		}
	}
}
