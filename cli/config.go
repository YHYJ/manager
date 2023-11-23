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
		installPath        = ""         // 定义在不同平台的安装路径
		installSourceTemp  = ""         // 定义在不同平台的Source安装方式的存储目录
		installReleaseTemp = ""         // 定义在不同平台的Release安装方式的存储目录
		goNames            = []string{} // 定义在不同平台可用的程序
		shellNames         = []string{} // 定义在不同平台可用的脚本
	)
	if general.Platform == "linux" {
		installPath = "/usr/local/bin"
		installSourceTemp = "/tmp/manager/source"
		installReleaseTemp = "/tmp/manager/release"
		goNames = []string{"checker", "clone-repos", "eniac", "kbdstage", "manager", "rolling", "scleaner", "skynet"}
		shellNames = []string{"collect-system", "configure-dtags", "py-virtualenv-tool", "save-docker-images", "sfm", "spacevim-update", "spider", "system-checkupdates", "trash-manager", "usb-manager"}
	} else if general.Platform == "darwin" {
		installPath = "/usr/local/bin"
		installSourceTemp = "/tmp/manager/source"
		installReleaseTemp = "/tmp/manager/release"
		goNames = []string{"clone-repos", "manager", "skynet"}
		shellNames = []string{"spacevim-update", "spider"}
	} else if general.Platform == "windows" {
		installPath = filepath.Join(general.GetVariable("ProgramFiles"), "Manager")
		installSourceTemp = filepath.Join(general.UserInfo.HomeDir, "AppData", "Local", "Temp", "manager", "source")
		installReleaseTemp = filepath.Join(general.UserInfo.HomeDir, "AppData", "Local", "Temp", "manager", "release")
		goNames = []string{"skynet"}
	}
	// 定义一个map[string]interface{}类型的变量并赋值
	exampleConf := map[string]interface{}{
		"variable": map[string]interface{}{ // 环境变量设置
			"http_proxy":  "http://<Domain or IP>:<Port>", // HTTP代理
			"https_proxy": "http://<Domain or IP>:<Port>", // HTTPS代理
		},
		"install": map[string]interface{}{
			"method":       "release",          // 安装方法，"release"或"source"代表安装预编译的二进制文件或自行从源码编译
			"path":         installPath,        // 安装路径
			"source_temp":  installSourceTemp,  // Source安装方式的基础存储目录
			"release_temp": installReleaseTemp, // Release安装方式的基础存储目录
			"go": map[string]interface{}{ // 基于go编写的程序的管理配置
				"names":                    goNames,                         // 可用的程序列表
				"release_api":              "https://api.github.com",        // Release安装的API地址
				"release_accept":           "application/vnd.github+json",   // Release安装的请求头参数
				"source_url":               "https://github.com",            // 首选Source安装源地址
				"source_api":               "https://api.github.com",        // 首选Source安装的API地址
				"source_username":          "YHYJ",                          // 首选Source安装的用户名
				"fallback_source_url":      "https://git.yj1516.top",        // 备用Source安装源地址
				"fallback_source_api":      "https://git.yj1516.top/api/v1", // 备用Source安装的API地址
				"fallback_source_username": "YJ",                            // 备用Source安装的用户名
				"unzip_path":               "flowout",                       // 解压缩释出文件夹
				"generate_path":            "build",                         // 编译结果存储文件夹
				"completion_dir": []string{ // zsh的自动补全文件夹
					filepath.Join(general.UserInfo.HomeDir, ".cache", "oh-my-zsh", "completions"),
					filepath.Join(general.UserInfo.HomeDir, ".oh-my-zsh", "cache", "completions"),
				},
			},
			"shell": map[string]interface{}{ // 基于shell编写的脚本的管理配置
				"source_api":               "https://api.github.com",              // 首选安装的API地址
				"source_raw":               "https://raw.githubusercontent.com",   // 首选安装的文件下载地址
				"source_username":          "YHYJ",                                // 首选安装的用户名
				"source_branch":            "ArchLinux",                           // 首选安装的分支名
				"fallback_source_api":      "https://git.yj1516.top/api/v1",       // 备用安装的API地址
				"fallback_source_raw":      "https://git.yj1516.top",              // 备用安装的文件下载地址
				"fallback_source_username": "YJ",                                  // 备用安装的用户名
				"fallback_source_branch":   "ArchLinux",                           // 备用安装的分支名
				"repo":                     "Program",                             // 存储脚本的仓库
				"dir":                      filepath.Join("System-Script", "app"), // 存储脚本的文件夹
				"names":                    shellNames,                            // 可用的脚本列表
			},
		},
	}
	// 检测配置文件是否存在
	if !general.FileExist(filePath) {
		return 0, fmt.Errorf("Open %s: no such file or directory", filePath)
	}
	// 检测配置文件是否是 toml 文件
	if !isTomlFile(filePath) {
		return 0, fmt.Errorf("Open %s: is not a toml file", filePath)
	}
	// 把 exampleConf 转换为 *toml.Tree 类型
	tree, err := toml.TreeFromMap(exampleConf)
	if err != nil {
		return 0, err
	}
	// 打开一个文件并获取 io.Writer 接口
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	return tree.WriteTo(file)
}
