/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 14:13:47

Description: 子命令 'setup' 的实现
*/

package cli

import (
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/yhyj/manager/general"
)

var (
	home       = general.GetVariable("HOME")
	hostname   = general.GetHostname()
	email      = "yj1516268@outlook.com"
	sep        = strings.Repeat(" ", 4)
	httpProxy  = "HTTP_PROXY=http://localhost:8080"
	httpsProxy = "HTTPS_PROXY=http://localhost:8080"
	noProxy    = "NO_PROXY=localhost,127.0.0.1,.example.com"

	// chezmoi 的依赖
	ChezmoiDependencies = "/usr/bin/chezmoi"                                        // 主程序
	ChezmoiConfigFile   = filepath.Join(home, ".config", "chezmoi", "chezmoi.toml") // 配置文件
	// chezmoi 配置
	chezmoiConfigFormat = "sourceDir = %s\n[git]\n%sautoCommit = %v\n%sautoPush = %v\n"
	chezmoiSourceDir    = `"~/Documents/Repos/System/Profile"`
	chezmoiAutoCommit   = false
	chezmoiAutoPush     = false
	ChezmoiConfig       = color.Sprintf(chezmoiConfigFormat, chezmoiSourceDir, sep, chezmoiAutoCommit, sep, chezmoiAutoPush)

	// cobra 的依赖
	CobraDependencies = filepath.Join(golangGOBIN, "cobra-cli") // 主程序
	CobraConfigFile   = filepath.Join(home, ".cobra.yaml")      // 配置文件
	// cobra 配置
	cobraConfigFormat = "author: %s <%s>\nlicense: %s\nuseViper: %v\n"
	cobraAuthor       = "YJ"
	cobraLicense      = "GPLv3"
	cobraUseViper     = false
	CobraConfig       = color.Sprintf(cobraConfigFormat, cobraAuthor, email, cobraLicense, cobraUseViper)

	// docker 的依赖
	DockerDependencies      = "/usr/bin/dockerd"                                   // 主程序
	DockerServiceConfigFile = "/etc/systemd/system/docker.service.d/override.conf" // 配置文件
	// docker 配置
	dockerServiceConfigFormat = "[Service]\nEnvironment=\"%s\"\nEnvironment=\"%s\"\nEnvironment=\"%s\"\nExecStart=\nExecStart=%s --data-root=%s -H fd://\n"
	dockerServiceDataRoot     = filepath.Join(home, "Documents", "Docker", "Root")
	DockerServiceConfig       = color.Sprintf(dockerServiceConfigFormat, httpProxy, httpsProxy, noProxy, DockerDependencies, dockerServiceDataRoot)

	// frpc 的依赖
	FrpcDependencies = "/usr/bin/frpc"                                    // 主程序
	FrpcConfigFile   = "/etc/systemd/system/frpc.service.d/override.conf" // 配置文件
	// frpc 配置
	frpcConfigFormat = "[Service]\nRestart=\nRestart=%s\n"
	frpcRestart      = "always"
	FrpcConfig       = color.Sprintf(frpcConfigFormat, frpcRestart)

	// git 的依赖
	GitDependencies = "/usr/bin/git"                    // 主程序
	GitConfigFile   = filepath.Join(home, ".gitconfig") // 配置文件
	// git 配置
	gitConfigFormat      = "[user]\n%sname = %s\n%semail = %s\n[core]\n%seditor = %s\n%sautocrlf = %s\n[merge]\n%stool = %s\n[color]\n%sui = %v\n[pull]\n%srebase = %v\n[filter \"lfs\"]\n%sclean = %s\n%ssmudge = %s\n%sprocess = %s\n%srequired = %v\n"
	gitCoreEditor        = "vim"
	gitCoreAutoCRLF      = "input"
	gitMergeTool         = "vimdiff"
	gitColorUI           = true
	gitPullRebase        = false
	gitFilterLfsClean    = "git-lfs clean -- %f"
	gitFilterLfsSmudge   = "git-lfs smudge -- %f"
	gitFilterLfsProcess  = "git-lfs filter-process"
	gitFilterLfsRequired = true
	GitConfig            = color.Sprintf(gitConfigFormat, sep, hostname, sep, email, sep, gitCoreEditor, sep, gitCoreAutoCRLF, sep, gitMergeTool, sep, gitColorUI, sep, gitPullRebase, sep, gitFilterLfsClean, sep, gitFilterLfsSmudge, sep, gitFilterLfsProcess, sep, gitFilterLfsRequired)

	// go 的依赖
	GolangDependencies = "/usr/bin/go"                               // 主程序
	GolangConfigFile   = filepath.Join(home, ".config", "go", "env") // 配置文件
	// go 配置
	golangConfigFormat = "GO111MODULE=%s\nGOBIN=%s\nGOPATH=%s\nGOCACHE=%s\nGOMODCACHE=%s\n"
	golangGO111MODULE  = "on"
	golangGOBIN        = filepath.Join(home, ".go", "bin")
	golangGOPATH       = filepath.Join(home, ".go")
	golangGOCACHE      = filepath.Join(home, ".cache", "go", "go-build")
	golangGOMODCACHE   = filepath.Join(home, ".cache", "go", "pkg", "mod")
	GolangConfig       = color.Sprintf(golangConfigFormat, golangGO111MODULE, golangGOBIN, golangGOPATH, golangGOCACHE, golangGOMODCACHE)

	// pip 的依赖
	PipDependencies = "/usr/bin/pip"                                    // 主程序
	PipConfigFile   = filepath.Join(home, ".config", "pip", "pip.conf") // 配置文件
	// pip 配置
	pipConfigFormat = "[global]\nindex-url = %s\ntrusted-host = %s\n"
	pipIndexUrl     = "https://mirrors.aliyun.com/pypi/simple"
	pipTrustedHost  = "mirrors.aliyun.com"
	PipConfig       = color.Sprintf(pipConfigFormat, pipIndexUrl, pipTrustedHost)

	// system-checkupdates Timer 和 Service 的依赖
	SystemCheckupdatesDependencies      = "/usr/local/bin/system-checkupdates"              // 主程序，需要版本 >= 3.0.0-20230313.1
	SystemCheckupdatesTimerConfigFile   = "/etc/systemd/system/system-checkupdates.timer"   // Timer 配置文件
	SystemCheckupdatesServiceConfigFile = "/etc/systemd/system/system-checkupdates.service" // Service 配置文件
	// system-checkupdates 配置 - Timer
	systemCheckupdatesTimerConfigFormat      = "[Unit]\nDescription=%s\n\n[Timer]\nOnBootSec=%s\nOnUnitInactiveSec=%s\nAccuracySec=%s\nPersistent=%v\n\n[Install]\nWantedBy=%s\n"
	systemcheckupdatesTimerDescription       = "Timer for system-checkupdates"
	systemcheckupdatesTimerOnBootSec         = "10min"
	systemcheckupdatesTimerOnUnitInactiveSec = "2h"
	systemcheckupdatesTimerAccuracySec       = "30min"
	systemcheckupdatesTimerPersistent        = true
	systemcheckupdatesTimerWantedBy          = "timers.target"
	SystemCheckupdatesTimerConfig            = color.Sprintf(systemCheckupdatesTimerConfigFormat, systemcheckupdatesTimerDescription, systemcheckupdatesTimerOnBootSec, systemcheckupdatesTimerOnUnitInactiveSec, systemcheckupdatesTimerAccuracySec, systemcheckupdatesTimerPersistent, systemcheckupdatesTimerWantedBy)
	// system-checkupdates 配置 - Service
	systemCheckupdatesServiceConfigFormat = "[Unit]\nDescription=%s\nAfter=%s\nWants=%s\n\n[Service]\nType=%s\nExecStart=%s\n"
	systemcheckupdatesServiceDescription  = "System checkupdates"
	systemcheckupdatesServiceAfter        = "network.target"
	systemcheckupdatesServiceWants        = "network.target"
	systemcheckupdatesServiceType         = "oneshot"
	systemcheckupdatesServiceExecStart    = "/usr/local/bin/system-checkupdates --check"
	SystemCheckupdatesServiceConfig       = color.Sprintf(systemCheckupdatesServiceConfigFormat, systemcheckupdatesServiceDescription, systemcheckupdatesServiceAfter, systemcheckupdatesServiceWants, systemcheckupdatesServiceType, systemcheckupdatesServiceExecStart)
)

// ProgramConfigurator 程序配置器
// 参数：
//   - flags: 系统信息各部分的开关
func ProgramConfigurator(flags map[string]bool) {
	// 预定义变量
	var (
		subjectName      string
		descriptorText   string
		subjectMinorName string
		writeMode        string = "t"
	)

	// 预定义输出格式
	var (
		subjectMinorNameFormat = "%*s%s %s\n"
		descriptorFormat       = "%*s%s %s: %s %s %s\n"
		configFileFormat       = "%*s%s %s: %s\n"
		errorFormat            = "%*s%s %s: %s\n\n"
		successFormat          = "%*s%s %s: %s\n\n"
	)

	// 配置 chezmoi
	if flags["chezmoiFlag"] {
		subjectName = "chezmoi"
		descriptorText = "configuration file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(ChezmoiConfigFile))
		if err := general.WriteFile(ChezmoiConfigFile, ChezmoiConfig, writeMode); err != nil {
			color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
	}

	// 配置 cobra
	if flags["cobraFlag"] {
		subjectName = "cobra-cli"
		descriptorText = "configuration file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(CobraConfigFile))
		if err := general.WriteFile(CobraConfigFile, CobraConfig, writeMode); err != nil {
			color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
	}

	// 配置 docker
	if flags["dockerFlag"] {
		subjectName = "docker"
		descriptorText = "daemon configuration"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(DockerServiceConfigFile))
		if err := general.WriteFile(DockerServiceConfigFile, DockerServiceConfig, writeMode); err != nil {
			color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
	}

	// 配置 frpc
	if flags["frpcFlag"] {
		subjectName = "frpc"
		descriptorText = "restart timing"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(FrpcConfigFile))
		if err := general.WriteFile(FrpcConfigFile, FrpcConfig, writeMode); err != nil {
			color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
	}

	// 配置 git
	if flags["gitFlag"] {
		subjectName = "git"
		descriptorText = "configuration file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(GitConfigFile))
		if err := general.WriteFile(GitConfigFile, GitConfig, writeMode); err != nil {
			color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
	}

	// 配置 golang
	if flags["goFlag"] {
		subjectName = "go"
		descriptorText = "environment file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(GolangConfigFile))
		if err := general.WriteFile(GolangConfigFile, GolangConfig, writeMode); err != nil {
			color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
	}

	// 配置 pip
	if flags["pipFlag"] {
		subjectName = "pip"
		descriptorText = "mirrors"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 1, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(PipConfigFile))
		if err := general.WriteFile(PipConfigFile, PipConfig, writeMode); err != nil {
			color.Printf(errorFormat, 1, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 1, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
	}

	// 配置 system-checkupdates
	if flags["systemcheckupdatesFlag"] {
		subjectName = "system-checkupdates"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		// system-checkupdates timer
		subjectMinorName = "system-checkupdates timer"
		descriptorText = "timer"
		color.Printf(subjectMinorNameFormat, 1, " ", general.SuccessText("-"), general.FgBlueText(subjectMinorName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(SystemCheckupdatesTimerConfigFile))
		if err := general.WriteFile(SystemCheckupdatesTimerConfigFile, SystemCheckupdatesTimerConfig, writeMode); err != nil {
			errorFormat = "%*s%s %s: %s\n"
			color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			successFormat = "%*s%s %s: %s\n"
			color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
		// system-checkupdates service
		subjectMinorName = "system-checkupdates service"
		descriptorText = "service"
		color.Printf(subjectMinorNameFormat, 1, " ", general.SuccessText("-"), general.FgBlueText(subjectMinorName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 1, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(SystemCheckupdatesServiceConfigFile))
		if err := general.WriteFile(SystemCheckupdatesServiceConfigFile, SystemCheckupdatesServiceConfig, writeMode); err != nil {
			color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
	}
}
