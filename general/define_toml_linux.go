//go:build linux

/*
File: define_toml_linux.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-11 15:15:12

Description: 操作 TOML 配置文件
*/

package general

import "path/filepath"

// 配置项
var (
	// 定义在不同平台的程序安装路径
	programPath = filepath.Join(Sep, "usr", "local", "bin")
	// 定义在不同平台的资源安装路径
	resourcesPath = filepath.Join(Sep, "usr", "local", "share")
	// 定义在不同平台的 Release 安装方式的存储目录
	releaseTemp = filepath.Join(Sep, "tmp", name, "release")
	// 定义在不同平台的 Source 安装方式的存储目录
	sourceTemp = filepath.Join(Sep, "tmp", name, "source")
	// 定义在不同平台的记账文件路径
	pocketPath = filepath.Join(Sep, "var", "local", "lib", name, "local")
	// 定义在不同平台可用的程序
	goNames = []string{
		name,
		"checker",
		"curator",
		"eniac",
		"kbdstage",
		"mqttc",
		"rolling",
		"scleaner",
		"skynet",
		"trash",
		"wocker",
	}
	// 定义在不同平台的自动补全文件路径（仅限 oh-my-zsh）
	goCompletionDir = []string{
		filepath.Join(UserInfo.HomeDir, ".cache", "oh-my-zsh", "completions"),
		filepath.Join(UserInfo.HomeDir, ".oh-my-zsh", "cache", "completions"),
	}
	// 定义在不同平台可用的脚本
	shellNames = []string{
		"configure-tags",
		"git-browser",
		"open",
		"py-virtualenv-tool",
		"spacevim-update",
		"spider",
		"syncer",
		"usb-manager",
	}
)
