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
	"path/filepath"
	"strings"
)

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
	giteaUrl       = "https://git.yj1516.top"
	giteaApi       = "https://git.yj1516.top/api/v1"
	giteaUsername  = "YJ"
)

// 配置
var appConfig = Config{
	Program: ProgramConfig{
		Method:      InstallMethod,
		ProgramPath: programPath,
		ReleaseTemp: releaseTemp,
		SourceTemp:  sourceTemp,
		PocketPath:  pocketPath,
		PocketFile:  pocketFile,
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
		},
	},
	Variable: VariableConfig{
		HTTPProxy:  HttpProxy,
		HTTPSProxy: HttpsProxy,
	},
}
