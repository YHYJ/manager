/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 14:13:47

Description: 子命令`setup`功能函数
*/

package cli

import (
	"fmt"
	"path/filepath"
	"strings"

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
	ChezmoiConfig       = fmt.Sprintf(chezmoiConfigFormat, chezmoiSourceDir, sep, chezmoiAutoCommit, sep, chezmoiAutoPush)
	ChezmoiConfigFile   = filepath.Join(home, ".config", "chezmoi", "chezmoi.toml")

	// cobra 的依赖项
	CobraDependencies = filepath.Join(goGOBIN, "cobra-cli")
	// cobra 配置
	cobraConfigFormat = "author: %s <%s>\nlicense: %s\nuseViper: %v\n"
	cobraAuthor       = "YJ"
	cobraLicense      = "GPLv3"
	cobraUseViper     = false
	CobraConfig       = fmt.Sprintf(cobraConfigFormat, cobraAuthor, email, cobraLicense, cobraUseViper)
	CobraConfigFile   = filepath.Join(home, ".cobra.yaml")

	// docker service 和 mirrors 的依赖项
	DockerDependencies = "/usr/bin/dockerd"
	// docker 配置 - docker service
	dockerServiceConfigFormat = "[Service]\nExecStart=\nExecStart=%s --data-root=%s -H fd://\n"
	dockerServiceExecStart    = "/usr/bin/dockerd"
	dockerServiceDataRoot     = filepath.Join(home, "Documents", "Docker", "Root")
	DockerServiceConfig       = fmt.Sprintf(dockerServiceConfigFormat, dockerServiceExecStart, dockerServiceDataRoot)
	DockerServiceConfigFile   = "/etc/systemd/system/docker.service.d/override.conf"
	// docker 配置 - docker mirrors
	dockerMirrorsConfigFormat    = "{\n%s\"registry-mirrors\": %s\n}\n"
	dockerMirrorsRegistryMirrors = []string{`"https://docker.mirrors.ustc.edu.cn"`}
	DockerMirrorsConfig          = fmt.Sprintf(dockerMirrorsConfigFormat, sep, dockerMirrorsRegistryMirrors)
	DockerMirrorsConfigFile      = "/etc/docker/daemon.json"

	// frpc 的依赖项
	FrpcDependencies = "/usr/bin/frpc"
	// frpc 配置
	frpcConfigFormat = "[Service]\nRestart=\nRestart=%s\n"
	frpcRestart      = "always"
	FrpcConfig       = fmt.Sprintf(frpcConfigFormat, frpcRestart)
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
	GitConfig            = fmt.Sprintf(gitConfigFormat, sep, hostname, sep, email, sep, gitCoreEditor, sep, gitCoreAutoCRLF, sep, gitMergeTool, sep, gitColorUI, sep, gitPullRebase, sep, gitFilterLfsClean, sep, gitFilterLfsSmudge, sep, gitFilterLfsProcess, sep, gitFilterLfsRequired)
	GitConfigFile        = filepath.Join(home, ".gitconfig")

	// go 的依赖项
	GoDependencies = "/usr/bin/go"
	// go 配置
	goConfigFormat = "GO111MODULE=%s\nGOBIN=%s\nGOPATH=%s\nGOCACHE=%s\nGOMODCACHE=%s\n"
	goGO111MODULE  = "on"
	goGOBIN        = filepath.Join(home, ".go", "bin")
	goGOPATH       = filepath.Join(home, ".go")
	goGOCACHE      = filepath.Join(home, ".cache", "go", "go-build")
	goGOMODCACHE   = filepath.Join(home, ".cache", "go", "pkg", "mod")
	GoConfig       = fmt.Sprintf(goConfigFormat, goGO111MODULE, goGOBIN, goGOPATH, goGOCACHE, goGOMODCACHE)
	GoConfigFile   = filepath.Join(home, ".config", "go", "env")

	// npm 的依赖项
	NpmDependencies = "/usr/bin/npm"
	// npm 配置
	npmConfigFormat = "registry=%s\n"
	npmRegistry     = "https://npmmirror.com"
	NpmConfig       = fmt.Sprintf(npmConfigFormat, npmRegistry)
	NpmConfigFile   = filepath.Join(home, ".npmrc")

	// pip 的依赖项
	PipDependencies = "/usr/bin/pip"
	// pip 配置
	pipConfigFormat = "[global]\nindex-url = %s\ntrusted-host = %s\n"
	pipIndexUrl     = "https://mirrors.aliyun.com/pypi/simple"
	pipTrustedHost  = "mirrors.aliyun.com"
	PipConfig       = fmt.Sprintf(pipConfigFormat, pipIndexUrl, pipTrustedHost)
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
	SystemCheckupdatesTimerConfig            = fmt.Sprintf(systemCheckupdatesTimerConfigFormat, systemcheckupdatesTimerDescription, systemcheckupdatesTimerOnBootSec, systemcheckupdatesTimerOnUnitInactiveSec, systemcheckupdatesTimerAccuracySec, systemcheckupdatesTimerPersistent, systemcheckupdatesTimerWantedBy)
	SystemCheckupdatesTimerConfigFile        = "/etc/systemd/system/system-checkupdates.timer"
	// system-checkupdates 配置 - system-checkupdates service
	systemCheckupdatesServiceConfigFormat = "[Unit]\nDescription=%s\nAfter=%s\nWants=%s\n\n[Service]\nType=%s\nExecStart=%s\n"
	systemcheckupdatesServiceDescription  = "System checkupdates"
	systemcheckupdatesServiceAfter        = "network.target"
	systemcheckupdatesServiceWants        = "network.target"
	systemcheckupdatesServiceType         = "oneshot"
	systemcheckupdatesServiceExecStart    = "/usr/local/bin/system-checkupdates --check"
	SystemCheckupdatesServiceConfig       = fmt.Sprintf(systemCheckupdatesServiceConfigFormat, systemcheckupdatesServiceDescription, systemcheckupdatesServiceAfter, systemcheckupdatesServiceWants, systemcheckupdatesServiceType, systemcheckupdatesServiceExecStart)
	SystemCheckupdatesServiceConfigFile   = "/etc/systemd/system/system-checkupdates.service"
)
