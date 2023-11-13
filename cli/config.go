/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:05:47

Description: 子命令`config`的实现
*/

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/yhyj/manager/general"
)

// isTomlFile 判断文件是不是 toml 文件
func isTomlFile(filePath string) bool {
	if strings.HasSuffix(filePath, ".toml") {
		return true
	}
	return false
}

// GetTomlConfig 读取 toml 配置文件
func GetTomlConfig(filePath string) (*toml.Tree, error) {
	if !general.FileExist(filePath) {
		return nil, fmt.Errorf("Open %s: no such file or directory", filePath)
	}
	if !isTomlFile(filePath) {
		return nil, fmt.Errorf("Open %s: is not a toml file", filePath)
	}
	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// WriteTomlConfig 写入 toml 配置文件
func WriteTomlConfig(filePath string) (int64, error) {
	// 根据系统不同决定某些参数
	var (
		installPath = ""         // 定义在不同平台的安装路径
		installTemp = ""         // 定义在不同平台的编译目录
		goNames     = []string{} // 定义在不同平台可用的程序
		shellNames  = []string{} // 定义在不同平台可用的脚本
	)
	if general.Platform == "linux" {
		installPath = "/usr/local/bin"
		installTemp = "/tmp/manager-build"
		goNames = []string{"checker", "clone-repos", "eniac", "kbdstage", "manager", "rolling", "scleaner", "skynet"}
		shellNames = []string{"collect-system", "configure-dtags", "py-virtualenv-tool", "save-docker-images", "sfm", "spacevim-update", "spider", "system-checkupdates", "trash-manager", "usb-manager"}
	} else if general.Platform == "darwin" {
		installPath = "/usr/local/bin"
		installTemp = "/tmp/manager-build"
		goNames = []string{"clone-repos", "manager", "skynet"}
		shellNames = []string{"spacevim-update", "spider"}
	} else if general.Platform == "windows" {
		installPath = filepath.Join(general.GetVariable("ProgramFiles"), "Manager")
		installTemp = filepath.Join(general.UserInfo.HomeDir, "AppData", "Local", "Temp")
		goNames = []string{"skynet"}
	}
	// 定义一个map[string]interface{}类型的变量并赋值
	exampleConf := map[string]interface{}{
		"variable": map[string]interface{}{
			"http_proxy":  "", // 例如"http://127.0.0.1:8080"
			"https_proxy": "", // 例如"http://127.0.0.1:8080"
		},
		"install": map[string]interface{}{
			"method": "release", // 安装方法，"release"或"source"代表安装预编译的二进制文件或自行从源码编译
			"path":   installPath,
			"temp":   installTemp,
			"go": map[string]interface{}{
				"generate_path":            "build",
				"source_url":               "https://git.yj1516.top",
				"source_username":          "YJ",
				"source_api":               "https://git.yj1516.top/api/v1",
				"fallback_source_url":      "https://github.com",
				"fallback_source_username": "YHYJ",
				"fallback_source_api":      "https://api.github.com",
				"names":                    goNames,
				"completion_dir":           []string{filepath.Join(general.UserInfo.HomeDir, ".cache", "oh-my-zsh", "completions"), filepath.Join(general.UserInfo.HomeDir, ".oh-my-zsh", "cache", "completions")},
			},
			"shell": map[string]interface{}{
				"source_url":               "https://git.yj1516.top",
				"source_username":          "YJ",
				"source_api":               "https://git.yj1516.top/api/v1",
				"source_branch":            "ArchLinux",
				"fallback_source_url":      "https://github.com",
				"fallback_source_username": "YHYJ",
				"fallback_source_api":      "https://api.github.com",
				"fallback_source_branch":   "ArchLinux",
				"repo":                     "Program",
				"dir":                      filepath.Join("System-Script", "app"),
				"names":                    shellNames,
			},
		},
	}
	if !general.FileExist(filePath) {
		return 0, fmt.Errorf("Open %s: no such file or directory", filePath)
	}
	if !isTomlFile(filePath) {
		return 0, fmt.Errorf("Open %s: is not a toml file", filePath)
	}
	// 把exampleConf转换为*toml.Tree
	tree, err := toml.TreeFromMap(exampleConf)
	if err != nil {
		return 0, err
	}
	// 打开一个文件并获取io.Writer接口
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	return tree.WriteTo(file)
}
