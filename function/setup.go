/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 14:13:47

Description: 子命令`setup`功能函数
*/

package function

var varHome = GetVariable("HOME")

// pip配置文件内容
var PipConfigFile = varHome + "/.config/pip/pip.conf"
var PipConfig = `[global]
index-url = http://mirrors.aliyun.com/pypi/simple
trusted-host = mirrors.aliyun.com
`

// npm配置文件内容
var NpmConfigFile = varHome + "/.npmrc"
var NpmConfig = `registry=https://registry.npm.taobao.org`

// docker配置文件内容
var DockerConfigFile = "/etc/systemd/system/docker.service.d/override.conf"
var DockerConfig = `[Service]
ExecStart=
ExecStart=/usr/bin/dockerd --data-root=` + varHome + `/Documents/Docker/Root -H fd://
`

// git配置文件内容
var GitConfigFile = varHome + "/.gitconfig"
var GitConfig = `[user]
	name = ` + GetHostname() + `
	email = yj1516268@outlook.com
[core]
	editor = /usr/bin/nvim
	autocrlf = input
[merge]
	tool = vimdiff
[color]
	ui = true
[pull]
	rebase = false
[filter "lfs"]
	clean = git-lfs clean -- %f
	smudge = git-lfs smudge -- %f
	process = git-lfs filter-process
	required = true
`
