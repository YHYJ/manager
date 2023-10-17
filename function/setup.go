/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 14:13:47

Description: 子命令`setup`功能函数
*/

package function

import (
	"fmt"
	"strings"
)

var (
	home     = GetVariable("HOME")
	hostname = GetHostname()
	email    = "yj1516268@outlook.com"
	sep      = strings.Repeat(" ", 4)

	// chezmoi
	chezmoiSourceDir    = `"~/Documents/Repos/System/Profile"`
	chezmoiAutoCommit   = false
	chezmoiAutoPush     = false
	chezmoiConfigFormat = "sourceDir = %s\n[git]\n%sautoCommit = %v\n%sautoPush = %v\n"
	ChezmoiConfig       = fmt.Sprintf(chezmoiConfigFormat, chezmoiSourceDir, sep, chezmoiAutoCommit, sep, chezmoiAutoPush)
	ChezmoiConfigFile   = home + "/.config/chezmoi/chezmoi.toml"

	// cobra配置
	cobraAuthor       = "YJ"
	cobraLicense      = "GPLv3"
	cobraUseViper     = false
	cobraConfigFormat = "author: %s <%s>\nlicense: %s\nuseViper: %v\n"
	CobraConfig       = fmt.Sprintf(cobraConfigFormat, cobraAuthor, email, cobraLicense, cobraUseViper)
	CobraConfigFile   = home + "/.cobra.yaml"

	// docker配置 - docker service
	dockerServiceExecStart    = "/usr/bin/dockerd"
	dockerServiceDataRoot     = home + "/Documents/Docker/Root"
	dockerServiceConfigFormat = "[Service]\nExecStart=\nExecStart=%s --data-root=%s -H fd://\n"
	DockerServiceConfig       = fmt.Sprintf(dockerServiceConfigFormat, dockerServiceExecStart, dockerServiceDataRoot)
	DockerServiceConfigFile   = "/etc/systemd/system/docker.service.d/override.conf"
	// docker配置 - docker mirrors
	dockerMirrorsRegistryMirrors = []string{`"https://docker.mirrors.ustc.edu.cn"`}
	dockerMirrorsConfigFormat    = "{\n%s\"registry-mirrors\": %s\n}\n"
	DockerMirrorsConfig          = fmt.Sprintf(dockerMirrorsConfigFormat, sep, dockerMirrorsRegistryMirrors)
	DockerMirrorsConfigFile      = "/etc/docker/daemon.json"

	// frpc配置
	frpcRestart      = "always"
	frpcConfigFormat = "[Service]\nRestart=\nRestart=%s\n"
	FrpcConfig       = fmt.Sprintf(frpcConfigFormat, frpcRestart)
	FrpcConfigFile   = "/etc/systemd/system/frpc.service.d/override.conf"

	// git配置
	gitCoreEditor        = "vim"
	gitCoreAutoCRLF      = "input"
	gitMergeTool         = "vimdiff"
	gitColorUI           = true
	gitPullRebase        = false
	gitFilterLfsClean    = "git-lfs clean -- %f"
	gitFilterLfsSmudge   = "git-lfs smudge -- %f"
	gitFilterLfsProcess  = "git-lfs filter-process"
	gitFilterLfsRequired = true
	gitConfigFormat      = "[user]\n%sname = %s\n%semail = %s\n[core]\n%seditor = %s\n%sautocrlf = %s\n[merge]\n%stool = %s\n[color]\n%sui = %v\n[pull]\n%srebase = %v\n[filter \"lfs\"]\n%sclean = %s\n%ssmudge = %s\n%sprocess = %s\n%srequired = %v\n"
	GitConfig            = fmt.Sprintf(gitConfigFormat, sep, hostname, sep, email, sep, gitCoreEditor, sep, gitCoreAutoCRLF, sep, gitMergeTool, sep, gitColorUI, sep, gitPullRebase, sep, gitFilterLfsClean, sep, gitFilterLfsSmudge, sep, gitFilterLfsProcess, sep, gitFilterLfsRequired)
	GitConfigFile        = home + "/.gitconfig"

	// golang配置
	goGO111MODULE  = "on"
	goGOBIN        = home + "/.go/bin"
	goGOPATH       = home + "/.go"
	goGOCACHE      = home + "/.cache/go/go-build"
	goGOMODCACHE   = home + "/.cache/go/pkg/mod"
	goConfigFormat = "GO111MODULE=%s\nGOBIN=%s\nGOPATH=%s\nGOCACHE=%s\nGOMODCACHE=%s\n"
	GoConfig       = fmt.Sprintf(goConfigFormat, goGO111MODULE, goGOBIN, goGOPATH, goGOCACHE, goGOMODCACHE)
	GoConfigFile   = home + "/.config/go/env"

	// npm配置
	npmRegistry     = "https://registry.npm.taobao.org"
	npmConfigFormat = "registry=%s\n"
	NpmConfig       = fmt.Sprintf(npmConfigFormat, npmRegistry)
	NpmConfigFile   = home + "/.npmrc"

	// pip配置
	pipIndexUrl     = "https://mirrors.aliyun.com/pypi/simple"
	pipTrustedHost  = "mirrors.aliyun.com"
	pipConfigFormat = "[global]\nindex-url = %s\ntrusted-host = %s\n"
	PipConfig       = fmt.Sprintf(pipConfigFormat, pipIndexUrl, pipTrustedHost)
	PipConfigFile   = home + "/.config/pip/pip.conf"

	// system-checkupdates timer和service的依赖项
	SystemCheckupdatesDependencies = "system-checkupdates >= 3.0.0-20230313.1"
	// system-checkupdates配置 - system-checkupdates timer
	systemCheckupdatesTimerConfigFormat      = "[Unit]\nDescription=%s\n\n[Timer]\nOnBootSec=%s\nOnUnitInactiveSec=%s\nAccuracySec=%s\nPersistent=%v\n\n[Install]\nWantedBy=%s\n"
	systemcheckupdatesTimerDescription       = "Timer for system-checkupdates"
	systemcheckupdatesTimerOnBootSec         = "10min"
	systemcheckupdatesTimerOnUnitInactiveSec = "2h"
	systemcheckupdatesTimerAccuracySec       = "30min"
	systemcheckupdatesTimerPersistent        = true
	systemcheckupdatesTimerWantedBy          = "timers.target"
	SystemCheckupdatesTimerConfig            = fmt.Sprintf(systemCheckupdatesTimerConfigFormat, systemcheckupdatesTimerDescription, systemcheckupdatesTimerOnBootSec, systemcheckupdatesTimerOnUnitInactiveSec, systemcheckupdatesTimerAccuracySec, systemcheckupdatesTimerPersistent, systemcheckupdatesTimerWantedBy)
	SystemCheckupdatesTimerConfigFile        = "/etc/systemd/system/system-checkupdates.timer"
	// system-checkupdates配置 - system-checkupdates service
	systemCheckupdatesServiceConfigFormat = "[Unit]\nDescription=%s\nAfter=%s\nWants=%s\n\n[Service]\nType=%s\nExecStart=%s\n"
	systemcheckupdatesServiceDescription  = "System checkupdates"
	systemcheckupdatesServiceAfter        = "network.target"
	systemcheckupdatesServiceWants        = "network.target"
	systemcheckupdatesServiceType         = "oneshot"
	systemcheckupdatesServiceExecStart    = "/usr/local/bin/system-checkupdates --check"
	SystemCheckupdatesServiceConfig       = fmt.Sprintf(systemCheckupdatesServiceConfigFormat, systemcheckupdatesServiceDescription, systemcheckupdatesServiceAfter, systemcheckupdatesServiceWants, systemcheckupdatesServiceType, systemcheckupdatesServiceExecStart)
	SystemCheckupdatesServiceConfigFile   = "/etc/systemd/system/system-checkupdates.service"
)
