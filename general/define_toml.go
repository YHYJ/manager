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
	"strings"
	"time"

	"github.com/gookit/color"
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
	// 打开配置文件
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// 写入注释
	manual := color.Sprintf("##\n## %s - %s\n## Generaled on %s\n##\n\n", Name, Version, time.Now().Format("2006-01-02 15:04:05"))
	n, err := file.WriteString(manual)
	if err != nil {
		return int64(n), err
	}

	// 创建编码器并设置顺序保留
	encoder := toml.NewEncoder(file)
	encoder.Order(toml.OrderPreserve)

	if err := encoder.Encode(appConfig); err != nil {
		return int64(n), err
	}

	stat, _ := file.Stat()
	return stat.Size(), nil
}
