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

// isTomlFile 检测文件是不是 toml 文件
//
// 参数：
//   - filePath: 待检测文件路径
//
// 返回：
//   - 是 toml 文件返回 true，否则返回 false
func isTomlFile(filePath string) bool {
	if strings.HasSuffix(filePath, ".toml") {
		return true
	}
	return false
}

// GetTomlConfig 读取 toml 配置文件
//
// 参数：
//   - filePath: toml 配置文件路径
//
// 返回：
//   - toml 配置树
//   - 错误信息
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
//
// 参数：
//   - filePath: toml 配置文件路径
//
// 返回：
//   - 写入的字节数
//   - 错误信息
func WriteTomlConfig(filePath string) (int64, error) {
	// 根据系统不同决定某些参数
	var (
		installProgramPath   = ""         // 定义在不同平台的程序安装路径
		installResourcesPath = ""         // 定义在不同平台的资源安装路径
		installSourceTemp    = ""         // 定义在不同平台的Source安装方式的存储目录
		installReleaseTemp   = ""         // 定义在不同平台的Release安装方式的存储目录
		goNames              = []string{} // 定义在不同平台可用的程序
		goCompletionDir      = []string{} // 定义在不同平台的自动补全文件路径（仅限oh-my-zsh）
		shellNames           = []string{} // 定义在不同平台可用的脚本
	)
	if general.Platform == "linux" {
		installProgramPath = "/usr/local/bin"
		installResourcesPath = "/usr/local/share"
		installSourceTemp = "/tmp/manager/source"
		installReleaseTemp = "/tmp/manager/release"
		goNames = []string{"checker", "curator", "eniac", "kbdstage", "manager", "rolling", "scleaner", "skynet", "trash"}
		goCompletionDir = []string{
			filepath.Join(general.UserInfo.HomeDir, ".cache", "oh-my-zsh", "completions"),
			filepath.Join(general.UserInfo.HomeDir, ".oh-my-zsh", "cache", "completions"),
		}
		shellNames = []string{
			"collect-system",
			"configure-dtags",
			"py-virtualenv-tool",
			"save-docker-images",
			"sfm",
			"spacevim-update",
			"spider",
			"system-checkupdates",
			"usb-manager",
		}
	} else if general.Platform == "darwin" {
		installProgramPath = "/usr/local/bin"
		installResourcesPath = "/usr/local/share"
		installSourceTemp = "/tmp/manager/source"
		installReleaseTemp = "/tmp/manager/release"
		goNames = []string{"curator", "manager", "skynet"}
		goCompletionDir = []string{
			filepath.Join(general.UserInfo.HomeDir, ".cache", "oh-my-zsh", "completions"),
			filepath.Join(general.UserInfo.HomeDir, ".oh-my-zsh", "cache", "completions"),
		}
		shellNames = []string{"spacevim-update"}
	} else if general.Platform == "windows" {
		installProgramPath = filepath.Join(general.GetVariable("ProgramFiles"), "Manager")
		installSourceTemp = filepath.Join(general.UserInfo.HomeDir, "AppData", "Local", "Temp", "manager", "source")
		installReleaseTemp = filepath.Join(general.UserInfo.HomeDir, "AppData", "Local", "Temp", "manager", "release")
		goNames = []string{"skynet"}
	}
	// 定义一个 map[string]interface{} 类型的变量并赋值
	exampleConf := map[string]interface{}{
		"variable": map[string]interface{}{ // 环境变量设置
			"http_proxy":  "", // HTTP 代理
			"https_proxy": "", // HTTPS 代理
		},
		"install": map[string]interface{}{
			"method":         "release",            // 安装方法，release 或 source 代表安装预编译的二进制文件或自行从源码编译
			"program_path":   installProgramPath,   // 程序安装路径
			"resources_path": installResourcesPath, // 资源安装路径
			"source_temp":    installSourceTemp,    // Source 安装方式的基础存储目录
			"release_temp":   installReleaseTemp,   // Release 安装方式的基础存储目录
			"go": map[string]interface{}{ // 基于 go 编写的程序的管理配置
				"names":           goNames,                         // 可用的程序列表
				"release_api":     "https://api.github.com",        // Release 安装源 API 地址
				"release_accept":  "application/vnd.github+json",   // Release 安装源请求头参数
				"generate_path":   "build",                         // Source 安装编译结果存储文件夹
				"github_url":      "https://github.com",            // Source 安装 - GitHub 安装源地址
				"github_api":      "https://api.github.com",        // Source 安装 - GitHub 安装源 API 地址
				"github_username": "YHYJ",                          // Source 安装 - GitHub 安装源用户名
				"gitea_url":       "https://git.yj1516.top",        // Source 安装 - Gitea 安装源地址
				"gitea_api":       "https://git.yj1516.top/api/v1", // Source 安装 - Gitea 安装源 API 地址
				"gitea_username":  "YJ",                            // Source 安装 - Gitea 安装源用户名
				"completion_dir":  goCompletionDir,                 // 自动补全文件夹
			},
			"shell": map[string]interface{}{ // 基于 shell 编写的脚本的管理配置
				"github_api":      "https://api.github.com",              // GitHub 安装源 API 地址
				"github_raw":      "https://raw.githubusercontent.com",   // GitHub 安装源文件下载地址
				"github_username": "YHYJ",                                // GitHub 安装源用户名
				"github_branch":   "ArchLinux",                           // GitHub 安装源分支名
				"gitea_api":       "https://git.yj1516.top/api/v1",       // Gitea 安装源 API 地址
				"gitea_raw":       "https://git.yj1516.top",              // Gitea 安装源文件下载地址
				"gitea_username":  "YJ",                                  // Gitea 安装源用户名
				"gitea_branch":    "ArchLinux",                           // Gitea 安装源分支名
				"repo":            "Program",                             // 存储脚本的仓库
				"dir":             filepath.Join("System-Script", "app"), // 脚本所在文件夹
				"names":           shellNames,                            // 可用的脚本列表
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
