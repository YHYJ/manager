/*
File: define_variable.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:01:45

Description: 操作变量
*/

package general

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/gookit/color"
)

// ---------- 代码变量

var (
	FgBlackText        = color.FgBlack.Render        // 前景色 - 黑色
	FgWhiteText        = color.FgWhite.Render        // 前景色 - 白色
	FgLightWhiteText   = color.FgLightWhite.Render   // 前景色 - 亮白色
	FgGrayText         = color.FgGray.Render         // 前景色 - 灰色
	FgRedText          = color.FgRed.Render          // 前景色 - 红色
	FgLightRedText     = color.FgLightRed.Render     // 前景色 - 亮红色
	FgGreenText        = color.FgGreen.Render        // 前景色 - 绿色
	FgLightGreenText   = color.FgLightGreen.Render   // 前景色 - 亮绿色
	FgYellowText       = color.FgYellow.Render       // 前景色 - 黄色
	FgLightYellowText  = color.FgLightYellow.Render  // 前景色 - 亮黄色
	FgBlueText         = color.FgBlue.Render         // 前景色 - 蓝色
	FgLightBlueText    = color.FgLightBlue.Render    // 前景色 - 亮蓝色
	FgMagentaText      = color.FgMagenta.Render      // 前景色 - 品红
	FgLightMagentaText = color.FgLightMagenta.Render // 前景色 - 亮品红
	FgCyanText         = color.FgCyan.Render         // 前景色 - 青色
	FgLightCyanText    = color.FgLightCyan.Render    // 前景色 - 亮青色

	BgBlackText        = color.BgBlack.Render        // 背景色 - 黑色
	BgWhiteText        = color.BgWhite.Render        // 背景色 - 白色
	BgLightWhiteText   = color.BgLightWhite.Render   // 背景色 - 亮白色
	BgGrayText         = color.BgGray.Render         // 背景色 - 灰色
	BgRedText          = color.BgRed.Render          // 背景色 - 红色
	BgLightRedText     = color.BgLightRed.Render     // 背景色 - 亮红色
	BgGreenText        = color.BgGreen.Render        // 背景色 - 绿色
	BgLightGreenText   = color.BgLightGreen.Render   // 背景色 - 亮绿色
	BgYellowText       = color.BgYellow.Render       // 背景色 - 黄色
	BgLightYellowText  = color.BgLightYellow.Render  // 背景色 - 亮黄色
	BgBlueText         = color.BgBlue.Render         // 背景色 - 蓝色
	BgLightBlueText    = color.BgLightBlue.Render    // 背景色 - 亮蓝色
	BgMagentaText      = color.BgMagenta.Render      // 背景色 - 品红
	BgLightMagentaText = color.BgLightMagenta.Render // 背景色 - 亮品红
	BgCyanText         = color.BgCyan.Render         // 背景色 - 青色
	BgLightCyanText    = color.BgLightCyan.Render    // 背景色 - 亮青色

	InfoText      = color.Info.Render      // Info 文本
	NoteText      = color.Note.Render      // Note 文本
	LightText     = color.Light.Render     // Light 文本
	ErrorText     = color.Error.Render     // Error 文本
	DangerText    = color.Danger.Render    // Danger 文本
	NoticeText    = color.Notice.Render    // Notice 文本
	SuccessText   = color.Success.Render   // Success 文本
	CommentText   = color.Comment.Render   // Comment 文本
	PrimaryText   = color.Primary.Render   // Primary 文本
	WarnText      = color.Warn.Render      // Warn 文本
	QuestionText  = color.Question.Render  // Question 文本
	SecondaryText = color.Secondary.Render // Secondary 文本
)

// 进度条变量
var ProgressParameters = map[string]string{
	"view": "0", // 是否显示进度条 0: 不显示 1: 显示
}

var (
	DownloadFlag = "📥"  // 运行状态符号 - 下载中
	LatestFlag   = "🎉"  // 运行状态符号 - 已是最新
	SuccessFlag  = "✅"  // 运行状态符号 - 成功
	WarningFlag  = "⚠️" // 运行状态符号 - 警告
	ErrorFlag    = "❌"  // 运行状态符号 - 失败
)

var (
	LatestVersionMessage       = "is already the latest version"                   // 输出文本 - 已安装的程序和脚本为最新版
	UnableToCompileMessage     = "Makefile or main.go file does not exist"         // 输出文本 - 缺失编译文件无法完成编译
	AcsInstallSuccessMessage   = "auto-completion script installed successfully"   // 输出文本 - 自动补全脚本安装成功
	AcsInstallFailedMessage    = "auto-completion script installation failed"      // 输出文本 - 自动补全脚本安装失败
	AcsUninstallSuccessMessage = "auto-completion script uninstalled successfully" // 输出文本 - 自动补全脚本卸载成功
	AcsUninstallFailedMessage  = "auto-completion script uninstallation failed"    // 输出文本 - 自动补全脚本卸载失败
)

var (
	GoLatestReleaseTagApiFormat      = "%s/repos/%s/%s/releases/latest" // API 和下载地址 - 请求远端仓库最新 Tag 的 API - Release
	GoLatestSourceTagApiFormat       = "%s/repos/%s/%s/tags"            // API 和下载地址 - 请求远端仓库最新 Tag 的 API - Source
	ShellLatestHashApiFormat         = "%s/repos/%s/%s/contents/%s/%s"  // API 和下载地址 - 请求远端仓库最新脚本的 Hash 值的 API
	ShellGithubBaseDownloadUrlFormat = "%s/%s/%s/%s"                    // API 和下载地址 - 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - GitHub 格式
	ShellGiteaBaseDownloadUrlFormat  = "%s/%s/%s/raw/branch/%s"         // API 和下载地址 - 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - Gitea 格式
)

var (
	CursorOnFlag   = "👉"  // 选择器符号 - 光标所在
	CursorOffFlag  = "  " // 选择器符号 - 光标不在
	SelectedFlag   = "•"  // 选择器符号 - 已选中
	UnselectedFlag = " "  // 选择器符号 - 已选中
	SelectAllFlag  = "⭐️" // 选择器符号 - 全选
)

var (
	Separator1st = "=" // 分隔符 - 1级
	Separator2st = "-" // 分隔符 - 2级
	Separator3st = "·" // 分隔符 - 3级
)

var (
	MultiSelectTips  = "Please select from the %s below (multi-select)\n"  // 提示词 - 多选
	SingleSelectTips = "Please select from the %s below (single-select)\n" // 提示词 - 单选
	QuietTips        = "Press '%s' to quit\n"                              // 提示词 - 退出
	SelectAllTips    = "Select All"                                        // 提示词 - 全选
	UninstallTips    = "Do you want to uninstall these software?"          // 提示词 - 卸载软件
)

// ---------- 环境变量

var Platform = runtime.GOOS                   // 操作系统
var Arch = runtime.GOARCH                     // 系统架构
var Sep = string(filepath.Separator)          // 路径分隔符
var UserInfo, _ = GetUserInfoByName(UserName) // 用户信息
// 用户名，当程序提权运行时，使用 SUDO_USER 变量获取提权前的用户名
var UserName = func() string {
	if GetVariable("SUDO_USER") != "" {
		return GetVariable("SUDO_USER")
	}
	return GetVariable("USER")
}()

// 用来处理不同系统之间的变量名差异
var platformChart = map[string]map[string]string{
	"windows": {
		"HOME":     "USERPROFILE",  // 用户主目录路径
		"USER":     "USERNAME",     // 当前登录用户名
		"SHELL":    "ComSpec",      // 默认shell或命令提示符路径
		"PWD":      "CD",           // 当前工作目录路径
		"HOSTNAME": "COMPUTERNAME", // 计算机主机名
	},
}

// GetVariable 获取环境变量
//
// 参数：
//   - key: 变量名
//
// 返回：
//   - 变量值
func GetVariable(key string) string {
	if innerMap, exists := platformChart[Platform]; exists {
		if _, variableExists := innerMap[key]; variableExists {
			key = platformChart[Platform][key]
		}
	}
	variable := os.Getenv(key)

	return variable
}

// GetHostname 获取系统 HOSTNAME
//
// 返回：
//   - HOSTNAME 或空字符串
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// SetVariable 设置环境变量
//
// 参数：
//   - key: 变量名
//   - value: 变量值
//
// 返回：
//   - 错误信息
func SetVariable(key, value string) error {
	return os.Setenv(key, value)
}

// GetUserInfoByName 根据用户名获取用户信息
//
// 参数：
//   - userName: 用户名
//
// 返回：
//   - 用户信息
//   - 错误信息
func GetUserInfoByName(userName string) (*user.User, error) {
	userInfo, err := user.Lookup(userName)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetUserInfoById 根据 ID 获取用户信息
//
// 参数：
//   - userId: 用户 ID
//
// 返回：
//   - 用户信息
//   - 错误信息
func GetUserInfoById(userId int) (*user.User, error) {
	userInfo, err := user.LookupId(strconv.Itoa(userId))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetCurrentUserInfo 获取当前用户信息
//
// 返回：
//   - 用户信息
//   - 错误信息
func GetCurrentUserInfo() (*user.User, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}
