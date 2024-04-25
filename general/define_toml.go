/*
File: define_toml.go
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

	"github.com/gookit/color"
	"github.com/pelletier/go-toml"
)

// 用于转换 Toml 配置树的结构体
type Config struct {
	Install  InstallConfig  `toml:"install"`
	Variable VariableConfig `toml:"variable"`
}
type InstallConfig struct {
	Method        string      `toml:"method"`
	ProgramPath   string      `toml:"program_path"`
	ReleaseTemp   string      `toml:"release_temp"`
	SourceTemp    string      `toml:"source_temp"`
	ResourcesPath string      `toml:"resources_path"`
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
	if Platform == "linux" {
		installProgramPath = "/usr/local/bin"
		installResourcesPath = "/usr/local/share"
		installSourceTemp = color.Sprintf("/tmp/%s/source", name)
		installReleaseTemp = color.Sprintf("/tmp/%s/release", name)
		goNames = []string{
			name,
			"checker",
			"curator",
			"eniac",
			"kbdstage",
			"rolling",
			"scleaner",
			"skynet",
			"trash",
		}
		goCompletionDir = []string{
			filepath.Join(UserInfo.HomeDir, ".cache", "oh-my-zsh", "completions"),
			filepath.Join(UserInfo.HomeDir, ".oh-my-zsh", "cache", "completions"),
		}
		shellNames = []string{
			"configure-dtags",
			"open",
			"open-remote-repository",
			"py-virtualenv-tool",
			"save-docker-images",
			"spacevim-update",
			"spider",
			"system-checkupdates",
			"usb-manager",
		}
	} else if Platform == "darwin" {
		installProgramPath = "/usr/local/bin"
		installResourcesPath = "/usr/local/share"
		installSourceTemp = color.Sprintf("/tmp/%s/source", name)
		installReleaseTemp = color.Sprintf("/tmp/%s/release", name)
		goNames = []string{
			name,
			"curator",
			"skynet",
		}
		goCompletionDir = []string{
			filepath.Join(UserInfo.HomeDir, ".cache", "oh-my-zsh", "completions"),
			filepath.Join(UserInfo.HomeDir, ".oh-my-zsh", "cache", "completions"),
		}
		shellNames = []string{
			"open-remote-repository",
			"spacevim-update",
			"spider",
		}
	} else if Platform == "windows" {
		installProgramPath = filepath.Join(GetVariable("ProgramFiles"), Name)
		installSourceTemp = filepath.Join(UserInfo.HomeDir, "AppData", "Local", "Temp", name, "source")
		installReleaseTemp = filepath.Join(UserInfo.HomeDir, "AppData", "Local", "Temp", name, "release")
		goNames = []string{
			name,
			"skynet",
		}
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
				"completion_dir":  goCompletionDir,                 // 自动补全文件夹
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
