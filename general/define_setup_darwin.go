//go:build darwin

/*
File: define_setup_darwin.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-07-30 14:53:39

Description: 供 'setup' 子命令使用的函数
*/

package general

import "path/filepath"

// ---------- 配置项

// Golang
var (
	// go 的依赖
	GolangDependencies = "go"                                                               // 主程序
	GolangConfigFile   = filepath.Join(home, "Library", "Application Support", "go", "env") // 配置文件
	// go 配置
	golangConfigFormat = "GO111MODULE=%s\nGOBIN=%s\nGOPATH=%s\nGOCACHE=%s\nGOMODCACHE=%s\n"
	golangGO111MODULE  = "on"
	golangGOPATH       = filepath.Join(home, ".go")
	golangGOCACHE      = filepath.Join(home, ".cache", "go", "go-build")
	golangGOMODCACHE   = filepath.Join(home, ".cache", "go", "pkg", "mod")
)

// Pip
var (
	// pip 的依赖
	PipDependencies = "pip"                                             // 主程序
	PipConfigFile   = filepath.Join(home, ".config", "pip", "pip.conf") // 配置文件
	// pip 配置
	pipConfigFormat = "[global]\nindex-url = %s\ntrusted-host = %s\n"
	pipIndexUrl     = "https://mirrors.aliyun.com/pypi/simple"
)
