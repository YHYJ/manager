//go:build windows

/*
File: define_toml_windows.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-11 15:15:12

Description: 操作 TOML 配置文件
*/

package general

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
)

// 用于转换 Toml 配置树的结构体
type Config struct {
	Program  ProgramConfig  `toml:"program"`
	Variable VariableConfig `toml:"variable"`
}
type ProgramConfig struct {
	Method        string      `toml:"method"`
	ProgramPath   string      `toml:"program_path"`
	ResourcesPath string      `toml:"resources_path"`
	ReleaseTemp   string      `toml:"release_temp"`
	SourceTemp    string      `toml:"source_temp"`
	PocketPath    string      `toml:"pocket_path"`
	PocketFile    string      `toml:"pocket_file"`
	Self          SelfConfig  `toml:"self"`
	Go            GoConfig    `toml:"go"`
	Shell         ShellConfig `toml:"shell"`
}
type VariableConfig struct {
	HTTPProxy  string `toml:"http_proxy"`
	HTTPSProxy string `toml:"https_proxy"`
}
type SelfConfig struct {
	Name           string   `toml:"name"`
	ReleaseApi     string   `toml:"release_api"`
	ReleaseAccept  string   `toml:"release_accept"`
	GeneratePath   string   `toml:"generate_path"`
	GithubUrl      string   `toml:"github_url"`
	GithubApi      string   `toml:"github_api"`
	GithubUsername string   `toml:"github_username"`
	GiteaUrl       string   `toml:"gitea_url"`
	GiteaApi       string   `toml:"gitea_api"`
	GiteaUsername  string   `toml:"gitea_username"`
	CompletionDir  []string `toml:"completion_dir"`
}
type GoConfig struct {
	Names          []string `toml:"names"`
	ReleaseApi     string   `toml:"release_api"`
	ReleaseAccept  string   `toml:"release_accept"`
	GeneratePath   string   `toml:"generate_path"`
	GithubUrl      string   `toml:"github_url"`
	GithubApi      string   `toml:"github_api"`
	GithubUsername string   `toml:"github_username"`
	GiteaUrl       string   `toml:"gitea_url"`
	GiteaApi       string   `toml:"gitea_api"`
	GiteaUsername  string   `toml:"gitea_username"`
	CompletionDir  []string `toml:"completion_dir"`
}
type ShellConfig struct {
	Names          []string `toml:"names"`
	Repo           string   `toml:"repo"`
	Dir            string   `toml:"dir"`
	GithubApi      string   `toml:"github_api"`
	GithubRaw      string   `toml:"github_raw"`
	GithubUsername string   `toml:"github_username"`
	GithubBranch   string   `toml:"github_branch"`
	GiteaApi       string   `toml:"gitea_api"`
	GiteaRaw       string   `toml:"gitea_raw"`
	GiteaUsername  string   `toml:"gitea_username"`
	GiteaBranch    string   `toml:"gitea_branch"`
}

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
	if !FileExist(filePath) {
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

// LoadConfigToStruct 将 Toml 配置树加载到结构体
//
// 参数：
//   - configTree: 解析 toml 配置文件得到的配置树
//
// 返回：
//   - 结构体
//   - 错误信息
func LoadConfigToStruct(configTree *toml.Tree) (*Config, error) {
	var config Config
	if err := configTree.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
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
	// 从常量读取该程序自身的名字
	name := strings.ToLower(Name)

	var (
		ProgramPath = ""         // 定义在不同平台的程序安装路径
		ReleaseTemp = ""         // 定义在不同平台的 Release 安装方式的存储目录
		SourceTemp  = ""         // 定义在不同平台的 Source 安装方式的存储目录
		PocketPath  = ""         // 定义在不同平台的记账文件路径
		goNames     = []string{} // 定义在不同平台可用的程序
	)
	if Platform == "windows" {
		ProgramPath = filepath.Join(GetVariable("ProgramFiles"), Name)
		ReleaseTemp = filepath.Join(UserInfo.HomeDir, "AppData", "Local", "Temp", name, "release")
		SourceTemp = filepath.Join(UserInfo.HomeDir, "AppData", "Local", "Temp", name, "source")
		PocketPath = filepath.Join(UserInfo.HomeDir, "AppData", "Local", "Temp", name, "local")
		goNames = []string{
			name,
			"skynet",
		}
	}

	// 定义一个 map[string]interface{} 类型的变量并赋值
	exampleConf := map[string]interface{}{
		"variable": map[string]interface{}{ // 环境变量设置
			"http_proxy":  HttpProxy,  // HTTP 代理
			"https_proxy": HttpsProxy, // HTTPS 代理
		},
		"program": map[string]interface{}{
			"method":       InstallMethod, // 安装方法，release 或 source 代表安装预编译的二进制文件或自行从源码编译
			"program_path": ProgramPath,   // 程序安装路径
			"release_temp": ReleaseTemp,   // Release 安装方式的基础存储目录
			"source_temp":  SourceTemp,    // Source 安装方式的基础存储目录
			"pocket_path":  PocketPath,    // 记账文件路径
			"pocket_file":  "files",       // 记账文件名
			"self": map[string]interface{}{ // 管理程序本身的配置
				"name":            strings.ToLower(Name),           // 管理程序名
				"release_api":     "https://api.github.com",        // Release 安装源 API 地址
				"release_accept":  "application/vnd.github+json",   // Release 安装源请求头参数
				"generate_path":   "build",                         // Source 安装编译结果存储文件夹
				"github_url":      "https://github.com",            // Source 安装 - GitHub 安装源地址
				"github_api":      "https://api.github.com",        // Source 安装 - GitHub 安装源 API 地址
				"github_username": "YHYJ",                          // Source 安装 - GitHub 安装源用户名
				"gitea_url":       "https://git.yj1516.top",        // Source 安装 - Gitea 安装源地址
				"gitea_api":       "https://git.yj1516.top/api/v1", // Source 安装 - Gitea 安装源 API 地址
				"gitea_username":  "YJ",                            // Source 安装 - Gitea 安装源用户名
			},
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
			},
		},
	}

	// 检测配置文件是否存在
	if !FileExist(filePath) {
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
