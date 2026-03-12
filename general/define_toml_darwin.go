//go:build darwin

/*
File: define_toml_darwin.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-11 15:15:12

Description: 操作 TOML 配置文件
*/

package general

import (
	"path/filepath"
	"strings"
)

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
		"curator",
		"eniac",
		"mqttc",
		"skynet",
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
		"py-virtualenv-tool",
		"spacevim-update",
		"spider",
		"trust-app",
	}
)

// 配置项
var (
	// 允许用户修改的配置项
	AllInstallMethod          = []string{"release", "source"}               // 所有安装方式，可选 source 或 release
	DefaultInstallMethodIndex = 0                                           // 默认安装方式的下标（从0开始）
	InstallMethod             = AllInstallMethod[DefaultInstallMethodIndex] // 默认安装方式
	HttpProxy                 = "http://127.0.0.1:8080"                     // 默认 HTTP 代理
	HttpsProxy                = HttpProxy                                   // 默认 HTTPS 代理，与 HTTP 代理一致
	// 使用默认值的配置项
	name           = strings.ToLower(Name)
	pocketFile     = "files"
	releaseApi     = "https://api.github.com"
	releaseAccept  = "application/vnd.github+json"
	generatePath   = "build"
	githubUrl      = "https://github.com"
	githubApi      = "https://api.github.com"
	githubUsername = "YHYJ"
	githubRaw      = "https://raw.githubusercontent.com"
	githubBranch   = "ArchLinux"
	giteaUrl       = "https://git.yj1516.top"
	giteaApi       = "https://git.yj1516.top/api/v1"
	giteaUsername  = "YJ"
	giteaRaw       = "https://git.yj1516.top"
	giteaBranch    = "ArchLinux"
	repo           = "Program"
	localF         = "System-Script"
	localC         = "app"
)

// 配置
var appConfig = Config{
	Program: ProgramConfig{
		Method:        InstallMethod,
		ProgramPath:   programPath,
		ResourcesPath: resourcesPath,
		ReleaseTemp:   releaseTemp,
		SourceTemp:    sourceTemp,
		PocketPath:    pocketPath,
		PocketFile:    pocketFile,
		Self: SelfConfig{
			Name:           name,
			ReleaseApi:     releaseApi,
			ReleaseAccept:  releaseAccept,
			GeneratePath:   generatePath,
			GithubUrl:      githubUrl,
			GithubApi:      githubApi,
			GithubUsername: githubUsername,
			GiteaUrl:       giteaUrl,
			GiteaApi:       giteaApi,
			GiteaUsername:  giteaUsername,
			CompletionDir:  goCompletionDir,
		},
		Go: GoConfig{
			Names:          goNames,
			ReleaseApi:     releaseApi,
			ReleaseAccept:  releaseAccept,
			GeneratePath:   generatePath,
			GithubUrl:      githubUrl,
			GithubApi:      githubApi,
			GithubUsername: githubUsername,
			GiteaUrl:       giteaUrl,
			GiteaApi:       giteaApi,
			GiteaUsername:  giteaUsername,
			CompletionDir:  goCompletionDir,
		},
		Shell: ShellConfig{
			Names:          shellNames,
			Repo:           repo,
			Dir:            filepath.Join(localF, localC),
			GithubApi:      githubApi,
			GithubRaw:      githubRaw,
			GithubUsername: githubUsername,
			GithubBranch:   githubBranch,
			GiteaApi:       giteaApi,
			GiteaRaw:       giteaRaw,
			GiteaUsername:  giteaUsername,
			GiteaBranch:    giteaBranch,
		},
	},
	Variable: VariableConfig{
		HTTPProxy:  HttpProxy,
		HTTPSProxy: HttpsProxy,
	},
}
