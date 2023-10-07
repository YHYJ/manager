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

	// pip配置文件内容
	PipConfig     = "[global]\nindex-url = http://mirrors.aliyun.com/pypi/simple\ntrusted-host = mirrors.aliyun.com"
	PipConfigFile = varHome + "/.config/pip/pip.conf"

	// npm配置文件内容
	NpmConfig     = "registry=https://registry.npm.taobao.org"
	NpmConfigFile = varHome + "/.npmrc"

	// docker配置文件内容
	DockerConfig     = "[Service]\nExecStart=\nExecStart=/usr/bin/dockerd --data-root=" + varHome + "/Documents/Docker/Root -H fd://"
	DockerConfigFile = "/etc/systemd/system/docker.service.d/override.conf"

	// frpc配置文件内容
	FrpcConfig     = "[Service]\nRestart=always\n"
	FrpcConfigFile = "/etc/systemd/system/frpc.service.d/override.conf"

	// git配置文件内容
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
	sep                  = "    "
	format               = "[user]\n%s%s\n%s%s\n[core]\n%s%s\n%s%s\n[merge]\n%s%s\n[color]\n%s%s\n[pull]\n%s%s\n[filter \"lfs\"]\n%s%s\n%s%s\n%s%s\n%s%s"
	GitConfig            = fmt.Sprintf(format, sep, gitUserName, sep, gitUserEmail, sep, gitCoreEditor, sep, gitCoreAutoCRLF, sep, gitMergeTool, sep, gitColorUI, sep, gitPullRebase, sep, gitFilterLfsClean, sep, gitFilterLfsSmudge, sep, gitFilterLfsProcess, sep, gitFilterLfsRequired)
	GitConfigFile        = varHome + "/.gitconfig"
)
