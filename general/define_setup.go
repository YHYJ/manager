/*
File: define_setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-07-30 14:53:39

Description: 供 'setup' 子命令使用的函数
*/

package general

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
)

// ---------- 预定义变量

// 预定义配置参数变量
var (
	hostname    = GetHostname()                      // 主机名
	userName, _ = GetUserName()                      // 用户真实名称
	home        = UserInfo.HomeDir                   // 用户家目录
	sep         = strings.Repeat(" ", 4)             // 分隔符
	goBin       = filepath.Join(home, ".go", "bin")  // `go install` 命令安装目录
	noProxy     = "localhost,127.0.0.1,.example.com" // 默认不代理的 URL
)

// 预定义配置项变量
var (
	subjectName      string       // 配置主题，应和程序/服务名称一致
	subjectMinorName string       // 配置子主题，在配置服务时识别是 timer 还是 service
	descriptorText   string       // 配置项描述
	writeMode        string = "t" // 配置文件写入模式，t 表示覆盖写入，a 表示追加写入
)

// 预定义输出格式
var (
	subjectMinorNameFormat = "%*s%s %s\n"
	descriptorFormat       = "%*s%s %s %s: %s %s %s\n"
	targetFileFormat       = "%*s%s %s %s: %s\n"
	askItemTitleFormat     = "%*s%s %s:\n"
	askItemsFormat         = "%*s%s "
	errorFormat            = "%*s%s %s: %s\n"
	statusFormat           = "%*s%s %s: %s\n"
	yesResultFormat        = "%*s%s %s\n"
	noResultFormat         = "%*s%s %s %s\n"
)

// ---------- 配置项

// Chezmoi
var (
	ChezmoiDependencies = "chezmoi"                                                 // 主程序
	ChezmoiConfigFile   = filepath.Join(home, ".config", "chezmoi", "chezmoi.toml") // 配置文件
	// chezmoi 配置
	chezmoiConfigFormat = "sourceDir = %s\n[git]\n%sautoCommit = %v\n%sautoPush = %v\n"
	chezmoiSourceDir    = `"~/Documents/Repos/System/Profile"`
	chezmoiAutoCommit   = "false"
	chezmoiAutoPush     = "false"
)

// Cobra
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

// Docker
var (
	// docker 的依赖
	DockerDependencies      = "dockerd"                                            // 主程序
	DockerServiceConfigFile = "/etc/systemd/system/docker.service.d/override.conf" // 配置文件
	// docker 配置
	dockerServiceConfigFormat = "[Service]\nEnvironment=\"%s\"\nEnvironment=\"%s\"\nEnvironment=\"%s\"\nExecStart=\nExecStart=%s --data-root=%s -H fd://\n"
	dockerServiceDataRoot     = filepath.Join(home, "Documents", "Docker", "Root")
)

// Frpc
var (
	// frpc 的依赖
	FrpcDependencies = "frpc"                                             // 主程序
	FrpcConfigFile   = "/etc/systemd/system/frpc.service.d/override.conf" // 配置文件
	// frpc 配置
	frpcConfigFormat = "[Service]\nRestart=\nRestart=%s\n"
	frpcRestart      = "always"
)

// Git
var (
	// git 的依赖
	GitDependencies = "git"                             // 主程序
	GitConfigFile   = filepath.Join(home, ".gitconfig") // 配置文件
	// git 配置
	gitConfigFormat      = "[user]\n%sname = %s\n%semail = %s\n[core]\n%seditor = %s\n%sautocrlf = %s\n[diff]\n%sexternal = %s\n[merge]\n%stool = %s\n[color]\n%sui = %v\n[pull]\n%srebase = %v\n[filter \"lfs\"]\n%sclean = %s\n%ssmudge = %s\n%sprocess = %s\n%srequired = %v\n"
	gitUserEmail         = "email@example.com"
	gitCoreEditor        = "vim"
	gitCoreAutoCRLF      = "input"
	gitDiffExternal      = "difft"
	gitMergeTool         = "vimdiff"
	gitColorUI           = "true"
	gitPullRebase        = "false"
	gitFilterLfsClean    = "git-lfs clean -- %f"
	gitFilterLfsSmudge   = "git-lfs smudge -- %f"
	gitFilterLfsProcess  = "git-lfs filter-process"
	gitFilterLfsRequired = "true"
)

// SystemCheckUpdates
var (
	// update-checker Timer 和 Service 的依赖
	UpdateCheckerDependencies = "checker"                                    // 主程序，需要版本 >= v0.7.0
	timerConfigFile           = "/etc/systemd/system/update-checker.timer"   // Timer 配置文件
	serviceConfigFile         = "/etc/systemd/system/update-checker.service" // Service 配置文件
	// update-checker 配置 - Timer
	timerConfigFormat      = "[Unit]\nDescription=%s\n\n[Timer]\nOnBootSec=%s\nOnUnitInactiveSec=%s\nAccuracySec=%s\nPersistent=%v\n\n[Install]\nWantedBy=%s\n"
	timerDescription       = "Timer for update-checker"
	timerOnBootSec         = "10min"
	timerOnUnitInactiveSec = "2h"
	timerAccuracySec       = "30min"
	timerPersistent        = "true"
	timerWantedBy          = "timers.target"
	// update-checker 配置 - Service
	serviceConfigFormat = "[Unit]\nDescription=%s\nAfter=%s\nWants=%s\n\n[Service]\nType=%s\nExecStart=%s\n"
	serviceDescription  = "Package update checker"
	serviceAfter        = "network.target"
	serviceWants        = "network.target"
	serviceType         = "oneshot"
	serviceExecStart    = "/usr/local/bin/checker update --file"
)

// SetupChezmoi 配置 Chezmoi
func SetupChezmoi() {
	// 提示
	subjectName = "chezmoi"
	descriptorText = "configuration file"
	color.Printf("%s %s\n", SuccessText("==>"), FgBlueText(subjectName))
	color.Printf(descriptorFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))

	// 检测
	if _, err := exec.LookPath(ChezmoiDependencies); err != nil {
		color.Printf(statusFormat, 2, " ", SuccessText("-"), LightText("Status"), NoticeText(color.Sprintf(InstallTips, subjectName)))
	} else {
		color.Printf(targetFileFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(ChezmoiConfigFile))
		// 创建配置文件
		if err := CreateFile(ChezmoiConfigFile); err != nil {
			color.Printf(errorFormat, 2, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			chezmoiSourceDir, _ = GetInput(QuestionText(color.Sprintf(InputTips, "sourceDir")), chezmoiSourceDir)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			chezmoiAutoCommit, _ = GetInput(QuestionText(color.Sprintf(InputTips, "autoCommit")), chezmoiAutoCommit)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			chezmoiAutoPush, _ = GetInput(QuestionText(color.Sprintf(InputTips, "autoPush")), chezmoiAutoPush)

			// 配置
			ChezmoiConfigContent := color.Sprintf(chezmoiConfigFormat, chezmoiSourceDir, sep, chezmoiAutoCommit, sep, chezmoiAutoPush)
			if err := WriteFile(ChezmoiConfigFile, ChezmoiConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 4, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 4, " ", SuccessText("-"), SuccessFlag)
			}
		}
	}
}

// SetupCobra 配置 Cobra
func SetupCobra() {
	// 提示
	subjectName = "cobra-cli"
	descriptorText = "configuration file"
	color.Printf("%s %s\n", SuccessText("==>"), FgBlueText(subjectName))
	color.Printf(descriptorFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))

	// 检测
	if _, err := exec.LookPath(CobraDependencies); err != nil {
		color.Printf(statusFormat, 2, " ", SuccessText("-"), LightText("Status"), NoticeText(color.Sprintf(InstallTips, subjectName)))
	} else {
		color.Printf(targetFileFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(CobraConfigFile))
		// 创建配置文件
		if err := CreateFile(CobraConfigFile); err != nil {
			color.Printf(errorFormat, 2, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			cobraAuthorName, _ = GetInput(QuestionText(color.Sprintf(InputTips, "authorName")), cobraAuthorName)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			cobraAuthorEmail, _ = GetInput(QuestionText(color.Sprintf(InputTips, "authorEmail")), cobraAuthorEmail)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			cobraLicense, _ = GetInput(QuestionText(color.Sprintf(InputTips, "license")), cobraLicense)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			cobraUseViper, _ = GetInput(QuestionText(color.Sprintf(InputTips, "useViper")), cobraUseViper)

			// 配置
			CobraConfigContent := color.Sprintf(cobraConfigFormat, cobraAuthorName, cobraAuthorEmail, cobraLicense, cobraUseViper)
			if err := WriteFile(CobraConfigFile, CobraConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 4, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 4, " ", SuccessText("-"), SuccessFlag)
			}
		}
	}
}

// SetupGit 配置 Git
func SetupGit() {
	// 提示
	subjectName = "git"
	descriptorText = "configuration file"
	color.Printf("%s %s\n", SuccessText("==>"), FgBlueText(subjectName))
	color.Printf(descriptorFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))

	// 检测
	if _, err := exec.LookPath(GitDependencies); err != nil {
		color.Printf(statusFormat, 2, " ", SuccessText("-"), LightText("Status"), NoticeText(color.Sprintf(InstallTips, subjectName)))
	} else {
		color.Printf(targetFileFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(GitConfigFile))
		// 创建配置文件
		if err := CreateFile(GitConfigFile); err != nil {
			color.Printf(errorFormat, 2, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			hostname, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[user].name")), hostname)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitUserEmail, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[user].email")), gitUserEmail)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitCoreEditor, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[core].editor")), gitCoreEditor)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitCoreAutoCRLF, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[core].autocrlf")), gitCoreAutoCRLF)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitDiffExternal, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[diff].external")), gitDiffExternal)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitMergeTool, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[merge].tool")), gitMergeTool)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitColorUI, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[color].ui")), gitColorUI)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitPullRebase, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[pull].rebase")), gitPullRebase)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitFilterLfsClean, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[filter \"lfs\"].clean")), gitFilterLfsClean)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitFilterLfsSmudge, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[filter \"lfs\"].smudge")), gitFilterLfsSmudge)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitFilterLfsProcess, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[filter \"lfs\"].process")), gitFilterLfsProcess)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			gitFilterLfsRequired, _ = GetInput(QuestionText(color.Sprintf(InputTips, "[filter \"lfs\"].required")), gitFilterLfsRequired)

			// 配置
			GitConfigContent := color.Sprintf(gitConfigFormat, sep, hostname, sep, gitUserEmail, sep, gitCoreEditor, sep, gitCoreAutoCRLF, sep, gitDiffExternal, sep, gitMergeTool, sep, gitColorUI, sep, gitPullRebase, sep, gitFilterLfsClean, sep, gitFilterLfsSmudge, sep, gitFilterLfsProcess, sep, gitFilterLfsRequired)
			if err := WriteFile(GitConfigFile, GitConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 4, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 4, " ", SuccessText("-"), SuccessFlag)
			}
		}
	}
}

// SetupGolang 配置 Golang
func SetupGolang() {
	// 提示
	subjectName = "go"
	descriptorText = "environment file"
	color.Printf("%s %s\n", SuccessText("==>"), FgBlueText(subjectName))
	color.Printf(descriptorFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))

	// 检测
	if _, err := exec.LookPath(GolangDependencies); err != nil {
		color.Printf(statusFormat, 2, " ", SuccessText("-"), LightText("Status"), NoticeText(color.Sprintf(InstallTips, subjectName)))
	} else {
		color.Printf(targetFileFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(GolangConfigFile))
		// 创建配置文件
		if err := CreateFile(GolangConfigFile); err != nil {
			color.Printf(errorFormat, 2, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			golangGO111MODULE, _ = GetInput(QuestionText(color.Sprintf(InputTips, "GO111MODULE")), golangGO111MODULE)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			golangGOPATH, _ = GetInput(QuestionText(color.Sprintf(InputTips, "GOPATH")), golangGOPATH)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			golangGOCACHE, _ = GetInput(QuestionText(color.Sprintf(InputTips, "GOCACHE")), golangGOCACHE)
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			golangGOMODCACHE, _ = GetInput(QuestionText(color.Sprintf(InputTips, "GOMODCACHE")), golangGOMODCACHE)

			// 配置
			GolangConfigContent := color.Sprintf(golangConfigFormat, golangGO111MODULE, goBin, golangGOPATH, golangGOCACHE, golangGOMODCACHE)
			if err := WriteFile(GolangConfigFile, GolangConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 4, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 4, " ", SuccessText("-"), SuccessFlag)
			}
		}
	}
}

// SetupPip 配置 Pip
func SetupPip() {
	// 提示
	subjectName = "pip"
	descriptorText = "mirrors"
	color.Printf("%s %s\n", SuccessText("==>"), FgBlueText(subjectName))
	color.Printf(descriptorFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Descriptor"), SecondaryText("Set up"), SecondaryText(subjectName), SecondaryText(descriptorText))

	// 检测
	if _, err := exec.LookPath(PipDependencies); err != nil {
		color.Printf(statusFormat, 2, " ", SuccessText("-"), LightText("Status"), NoticeText(color.Sprintf(InstallTips, subjectName)))
	} else {
		color.Printf(targetFileFormat, 2, " ", SuccessText("-"), InfoText("INFO:"), LightText("Target file"), CommentText(PipConfigFile))
		// 创建配置文件
		if err := CreateFile(PipConfigFile); err != nil {
			color.Printf(errorFormat, 2, " ", SuccessText("-"), LightText("Error"), DangerText(err))
		} else {
			// 交互
			color.Printf(askItemTitleFormat, 2, " ", SuccessText("-"), LightText("Config"))
			color.Printf(askItemsFormat, 4, " ", SuccessText("-"))
			pipIndexUrl, _ = GetInput(QuestionText(color.Sprintf(InputTips, "index-url")), pipIndexUrl)

			// 需要获取交互结果的配置项
			pipTrustedHost, _ := GetUrlHost(pipIndexUrl)

			// 配置
			PipConfigContent := color.Sprintf(pipConfigFormat, pipIndexUrl, pipTrustedHost)
			if err := WriteFile(PipConfigFile, PipConfigContent, writeMode); err != nil {
				color.Printf(noResultFormat, 4, " ", SuccessText("-"), ErrorFlag, DangerText(err))
			} else {
				color.Printf(yesResultFormat, 4, " ", SuccessText("-"), SuccessFlag)
			}
		}
	}
}
