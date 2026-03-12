//go:build windows

/*
File: define_toml_windows.go
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
	programPath = filepath.Join(GetVariable("ProgramFiles"), Name)
	// 定义在不同平台的 Release 安装方式的存储目录
	releaseTemp = filepath.Join(UserInfo.HomeDir, "AppData", "Local", "Temp", name, "release")
	// 定义在不同平台的 Source 安装方式的存储目录
	sourceTemp = filepath.Join(UserInfo.HomeDir, "AppData", "Local", "Temp", name, "source")
	// 定义在不同平台的记账文件路径
	pocketPath = filepath.Join(UserInfo.HomeDir, "AppData", "Local", "Temp", name, "local")
	// 定义在不同平台可用的程序
	goNames = []string{
		name,
		"mqttc",
		"skynet",
		"wocker",
	}
)
