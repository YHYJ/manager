/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 14:13:47

Description: 子命令`setup`功能函数
*/

package function

import "fmt"

var (
	varHome = GetVariable("HOME")
	sep     = "    "

	// chezmoi
	chezmoiSourceDir    = `"~/Documents/Repos/System/Profile"`
	chezmoiAutoCommit   = false
	chezmoiAutoPush     = false
	chezmoiConfigFormat = "sourceDir = %s\n[git]\n%sautoCommit = %v\n%sautoPush = %v\n"
	ChezmoiConfig       = fmt.Sprintf(chezmoiConfigFormat, chezmoiSourceDir, sep, chezmoiAutoCommit, sep, chezmoiAutoPush)
	ChezmoiConfigFile   = varHome + "/.config/chezmoi/chezmoi.toml"

	// cobra配置
	cobraAuthor       = "YJ"
	cobraEmail        = "yj1516268@outlook.com"
	cobraLicense      = "GPLv3"
	cobraUseViper     = false
	cobraConfigFormat = "author: %s <%s>\nlicense: %s\nuseViper: %v\n"
	CobraConfig       = fmt.Sprintf(cobraConfigFormat, cobraAuthor, cobraEmail, cobraLicense, cobraUseViper)
	CobraConfigFile   = varHome + "/.cobra.yaml"

	// docker配置
	// docker service
	DockerServiceConfig     = "[Service]\nExecStart=\nExecStart=/usr/bin/dockerd --data-root=" + varHome + "/Documents/Docker/Root -H fd://\n"
	DockerServiceConfigFile = "/etc/systemd/system/docker.service.d/override.conf"
	// docker mirrors
	dockerMirrorsRegistryMirrors = []string{`"https://docker.mirrors.ustc.edu.cn"`}
	dockerMirrorsConfigFormat    = "{\n%s\"registry-mirrors\": %s\n}\n"
	DockerMirrorsConfig          = fmt.Sprintf(dockerMirrorsConfigFormat, sep, dockerMirrorsRegistryMirrors)
	DockerMirrorsConfigFile      = "/etc/docker/daemon.json"

	// frpc配置
	FrpcConfig     = "[Service]\nRestart=always\n"
	FrpcConfigFile = "/etc/systemd/system/frpc.service.d/override.conf"

	// git配置
	gitUserName          = "name = " + GetHostname()
	gitUserEmail         = "email = yj1516268@outlook.com"
	gitCoreEditor        = "editor = vim"
	gitCoreAutoCRLF      = "autocrlf = input"
	gitMergeTool         = "tool = vimdiff"
	gitColorUI           = "ui = true"
	gitPullRebase        = "rebase = false"
	gitFilterLfsClean    = "clean = git-lfs clean -- %f"
	gitFilterLfsSmudge   = "smudge = git-lfs smudge -- %f"
	gitFilterLfsProcess  = "process = git-lfs filter-process"
	gitFilterLfsRequired = "required = true"
	gitConfigFormat      = "[user]\n%s%s\n%s%s\n[core]\n%s%s\n%s%s\n[merge]\n%s%s\n[color]\n%s%s\n[pull]\n%s%s\n[filter \"lfs\"]\n%s%s\n%s%s\n%s%s\n%s%s\n"
	GitConfig            = fmt.Sprintf(gitConfigFormat, sep, gitUserName, sep, gitUserEmail, sep, gitCoreEditor, sep, gitCoreAutoCRLF, sep, gitMergeTool, sep, gitColorUI, sep, gitPullRebase, sep, gitFilterLfsClean, sep, gitFilterLfsSmudge, sep, gitFilterLfsProcess, sep, gitFilterLfsRequired)
	GitConfigFile        = varHome + "/.gitconfig"

	// golang配置
	GoConfig     = fmt.Sprintf("GO111MODULE=on\nGOBIN=%s/.go/bin\nGOPATH=%s/.go\nGOCACHE=%s/.cache/go/go-build\nGOMODCACHE=%s/.cache/go/pkg/mod\n", varHome, varHome, varHome, varHome)
	GoConfigFile = varHome + "/.config/go/env"

	// npm配置
	NpmConfig     = "registry=https://registry.npm.taobao.org\n"
	NpmConfigFile = varHome + "/.npmrc"

	// pip配置
	PipConfig     = "[global]\nindex-url = http://mirrors.aliyun.com/pypi/simple\ntrusted-host = mirrors.aliyun.com\n"
	PipConfigFile = varHome + "/.config/pip/pip.conf"
)
