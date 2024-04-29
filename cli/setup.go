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
	home     = general.GetVariable("HOME")
	hostname = general.GetHostname()
	email    = "yj1516268@outlook.com"
	sep      = strings.Repeat(" ", 4)

	// chezmoi 的依赖项
	ChezmoiDependencies = "/usr/bin/chezmoi"
	// chezmoi 配置
	chezmoiConfigFormat = "sourceDir = %s\n[git]\n%sautoCommit = %v\n%sautoPush = %v\n"
	chezmoiSourceDir    = `"~/Documents/Repos/System/Profile"`
	chezmoiAutoCommit   = false
	chezmoiAutoPush     = false
	ChezmoiConfig       = color.Sprintf(chezmoiConfigFormat, chezmoiSourceDir, sep, chezmoiAutoCommit, sep, chezmoiAutoPush)
	ChezmoiConfigFile   = filepath.Join(home, ".config", "chezmoi", "chezmoi.toml")

	// cobra 的依赖项
	CobraDependencies = filepath.Join(golangGOBIN, "cobra-cli")
	// cobra 配置
	cobraConfigFormat = "author: %s <%s>\nlicense: %s\nuseViper: %v\n"
	cobraAuthor       = "YJ"
	cobraLicense      = "GPLv3"
	cobraUseViper     = false
	CobraConfig       = color.Sprintf(cobraConfigFormat, cobraAuthor, email, cobraLicense, cobraUseViper)
	CobraConfigFile   = filepath.Join(home, ".cobra.yaml")

	// docker service 和 mirrors 的依赖项
	DockerDependencies = "/usr/bin/dockerd"
	// docker 配置 - docker service
	dockerServiceConfigFormat = "[Service]\nExecStart=\nExecStart=%s --data-root=%s -H fd://\n"
	dockerServiceExecStart    = "/usr/bin/dockerd"
	dockerServiceDataRoot     = filepath.Join(home, "Documents", "Docker", "Root")
	DockerServiceConfig       = color.Sprintf(dockerServiceConfigFormat, dockerServiceExecStart, dockerServiceDataRoot)
	DockerServiceConfigFile   = "/etc/systemd/system/docker.service.d/override.conf"
	// docker 配置 - docker mirrors
	dockerMirrorsConfigFormat    = "{\n%s\"registry-mirrors\": %s\n}\n"
	dockerMirrorsRegistryMirrors = []string{`"https://docker.mirrors.ustc.edu.cn"`}
	DockerMirrorsConfig          = color.Sprintf(dockerMirrorsConfigFormat, sep, dockerMirrorsRegistryMirrors)
	DockerMirrorsConfigFile      = "/etc/docker/daemon.json"

	// frpc 的依赖项
	FrpcDependencies = "/usr/bin/frpc"
	// frpc 配置
	frpcConfigFormat = "[Service]\nRestart=\nRestart=%s\n"
	frpcRestart      = "always"
	FrpcConfig       = color.Sprintf(frpcConfigFormat, frpcRestart)
	FrpcConfigFile   = "/etc/systemd/system/frpc.service.d/override.conf"

	// git 的依赖项
	GitDependencies = "/usr/bin/git"
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
	GitConfigFile        = filepath.Join(home, ".gitconfig")

	// go 的依赖项
	GolangDependencies = "/usr/bin/go"
	// go 配置
	golangConfigFormat = "GO111MODULE=%s\nGOBIN=%s\nGOPATH=%s\nGOCACHE=%s\nGOMODCACHE=%s\n"
	golangGO111MODULE  = "on"
	golangGOBIN        = filepath.Join(home, ".go", "bin")
	golangGOPATH       = filepath.Join(home, ".go")
	golangGOCACHE      = filepath.Join(home, ".cache", "go", "go-build")
	golangGOMODCACHE   = filepath.Join(home, ".cache", "go", "pkg", "mod")
	GolangConfig       = color.Sprintf(golangConfigFormat, golangGO111MODULE, golangGOBIN, golangGOPATH, golangGOCACHE, golangGOMODCACHE)
	GolangConfigFile   = filepath.Join(home, ".config", "go", "env")

	// pip 的依赖项
	PipDependencies = "/usr/bin/pip"
	// pip 配置
	pipConfigFormat = "[global]\nindex-url = %s\ntrusted-host = %s\n"
	pipIndexUrl     = "https://mirrors.aliyun.com/pypi/simple"
	pipTrustedHost  = "mirrors.aliyun.com"
	PipConfig       = color.Sprintf(pipConfigFormat, pipIndexUrl, pipTrustedHost)
	PipConfigFile   = filepath.Join(home, ".config", "pip", "pip.conf")

	// system-checkupdates timer 和 service 的依赖项
	SystemCheckupdatesDependencies = "/usr/local/bin/system-checkupdates" // >= 3.0.0-20230313.1
	// system-checkupdates 配置 - system-checkupdates timer
	systemCheckupdatesTimerConfigFormat      = "[Unit]\nDescription=%s\n\n[Timer]\nOnBootSec=%s\nOnUnitInactiveSec=%s\nAccuracySec=%s\nPersistent=%v\n\n[Install]\nWantedBy=%s\n"
	systemcheckupdatesTimerDescription       = "Timer for system-checkupdates"
	systemcheckupdatesTimerOnBootSec         = "10min"
	systemcheckupdatesTimerOnUnitInactiveSec = "2h"
	systemcheckupdatesTimerAccuracySec       = "30min"
	systemcheckupdatesTimerPersistent        = true
	systemcheckupdatesTimerWantedBy          = "timers.target"
	SystemCheckupdatesTimerConfig            = color.Sprintf(systemCheckupdatesTimerConfigFormat, systemcheckupdatesTimerDescription, systemcheckupdatesTimerOnBootSec, systemcheckupdatesTimerOnUnitInactiveSec, systemcheckupdatesTimerAccuracySec, systemcheckupdatesTimerPersistent, systemcheckupdatesTimerWantedBy)
	SystemCheckupdatesTimerConfigFile        = "/etc/systemd/system/system-checkupdates.timer"
	// system-checkupdates 配置 - system-checkupdates service
	systemCheckupdatesServiceConfigFormat = "[Unit]\nDescription=%s\nAfter=%s\nWants=%s\n\n[Service]\nType=%s\nExecStart=%s\n"
	systemcheckupdatesServiceDescription  = "System checkupdates"
	systemcheckupdatesServiceAfter        = "network.target"
	systemcheckupdatesServiceWants        = "network.target"
	systemcheckupdatesServiceType         = "oneshot"
	systemcheckupdatesServiceExecStart    = "/usr/local/bin/system-checkupdates --check"
	SystemCheckupdatesServiceConfig       = color.Sprintf(systemCheckupdatesServiceConfigFormat, systemcheckupdatesServiceDescription, systemcheckupdatesServiceAfter, systemcheckupdatesServiceWants, systemcheckupdatesServiceType, systemcheckupdatesServiceExecStart)
	SystemCheckupdatesServiceConfigFile   = "/etc/systemd/system/system-checkupdates.service"
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

	// 定义输出格式
	subjectMinorNameFormat := "%*s%s %s\n"
	descriptorFormat := "%*s%s %s: %s %s %s\n"
	configFileFormat := "%*s%s %s: %s\n"
	errorFormat := "%*s%s %s: %s\n\n"
	successFormat := "%*s%s %s: %s\n\n"

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
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		// docker service
		subjectMinorName = "docker service"
		descriptorText = "root directory"
		color.Printf(subjectMinorNameFormat, 1, " ", general.SuccessText("-"), general.FgBlueText(subjectMinorName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 2, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(DockerServiceConfigFile))
		if err := general.WriteFile(DockerServiceConfigFile, DockerServiceConfig, writeMode); err != nil {
			errorFormat = "%*s%s %s: %s\n"
			color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			successFormat = "%*s%s %s: %s\n"
			color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
		}
		// docker mirrors
		subjectMinorName = "docker mirrors"
		descriptorText = "registry mirrors"
		color.Printf(subjectMinorNameFormat, 1, " ", general.SuccessText("-"), general.FgBlueText(subjectMinorName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.LightText("Descriptor"), general.CommentText("Set up"), general.CommentText(subjectName), general.CommentText(descriptorText))
		color.Printf(configFileFormat, 2, " ", general.SuccessText("-"), general.LightText("Configuration file"), general.CommentText(DockerMirrorsConfigFile))
		if err := general.WriteFile(DockerMirrorsConfigFile, DockerMirrorsConfig, writeMode); err != nil {
			color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err.Error()))
		} else {
			color.Printf(successFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.BgYellowText("Setup completed"))
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
