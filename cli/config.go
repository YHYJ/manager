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
	// 定义一个map[string]interface{}类型的变量并赋值
	exampleConf := map[string]interface{}{
		"variable": map[string]interface{}{
			"http_proxy":  "",
			"https_proxy": "",
		},
		"install": map[string]interface{}{
			"path": "/usr/local/bin",
			"temp": "/tmp/manager-build",
			"go": map[string]interface{}{
				"generate_path":            "build",
				"source_url":               "https://git.yj1516.top",
				"source_username":          "YJ",
				"source_api":               "https://git.yj1516.top/api/v1",
				"fallback_source_url":      "https://github.com",
				"fallback_source_username": "YHYJ",
				"fallback_source_api":      "https://api.github.com",
				"names":                    []string{"checker", "clone-repos", "eniac", "kbdstage", "manager", "rolling", "scleaner", "skynet"},
				"completion_dir":           []string{general.UserInfo.HomeDir + "/.cache/oh-my-zsh/completions", general.UserInfo.HomeDir + "/.oh-my-zsh/cache/completions"},
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
				"dir":                      "System-Script/app",
				"names":                    []string{"collect-system", "configure-dtags", "py-virtualenv-tool", "save-docker-images", "sfm", "spacevim-update", "spider", "system-checkupdates", "trash-manager", "usb-manager"},
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
