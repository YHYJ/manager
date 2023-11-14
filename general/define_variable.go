/*
File: define_variable.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:01:45

Description: 操作变量（包括代码变量和环境变量）
*/

package general

import (
	"os"
	"os/user"
	"runtime"
	"strconv"
)

// ---------- 代码变量

// 输出格式变量
var (
	// 常规输出格式
	Regelar1PFormat = "%s\n"   // 常规输出格式·1部分: <输出内容>
	Regelar2PFormat = "%s%s\n" // 常规输出格式·2部分: <输出内容1><输出内容2>
	// 标题部分输出格式
	TitleH1Format = "\n\x1b[36;3m%s\x1b[0m\n\n" // H1级别标题: <标题文字>
	// 分隔线输出格式
	LineHiddenFormat = "\x1b[30m%s\x1b[0m\n"   // 隐性分隔线: <分隔线>
	LineShownFormat  = "\x1b[30;1m%s\x1b[0m\n" // 显性分隔线: <分隔线>
	// Slice输出格式
	SliceTraverse2PFormat                = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m\n"                                             // 切片遍历·2部分: <指示符><前缀><元素>
	SliceTraverse2PSuffixFormat          = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m%s%s\n"                                         // 带后缀的切片遍历·2部分: <指示符><前缀><元素><分隔符><后缀>
	SliceTraverse2PSuffixNoNewLineFormat = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m%s%s"                                           // 带后缀的切片遍历·2部分·不换行: <指示符><前缀><元素><分隔符><后缀>
	SliceTraverse3PSuffixFormat          = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m%s\x1b[33;1m%s\x1b[0m%s%s\n"                    // 带后缀的切片遍历·3部分: <指示符><前缀><元素><分隔符><后缀1><分隔符><后缀2>
	SliceTraverse4PSuffixFormat          = "\x1b[32;1m%s\x1b[0m%s\x1b[34m%s\x1b[0m%s\x1b[33m%s\x1b[0m%s\x1b[35;1m%s\x1b[0m%s%s\n" // 带后缀的切片遍历·4部分: <指示符><前缀><元素><分隔符><后缀1><分隔符><后缀2><分隔符><后缀3>
	// 成功信息输出格式
	SuccessFormat       = "\x1b[32;1m%s\x1b[0m\n"     // 成功信息: <成功信息>
	SuccessSuffixFormat = "\x1b[32;1m%s\x1b[0m%s%s\n" // 带后缀的成功信息: <成功信息><分隔符><后缀信息>
	// 提示信息输出格式
	InfoFormat             = "\x1b[33;1m%s\x1b[0m\n"         // 提示信息: <提示信息>
	InfoPrefixFormat       = "%s%s\x1b[33;1m%s\x1b[0m\n"     // 带前缀的提示信息: <前缀信息><分隔符><提示信息>
	InfoSuffixFormat       = "\x1b[33;1m%s\x1b[0m%s%s\n"     // 带后缀的提示信息: <提示信息><分隔符><后缀信息>
	InfoPrefixSuffixFormat = "%s%s\x1b[33;1m%s\x1b[0m%s%s\n" // 带前后缀的提示信息: <前缀信息><分隔符><提示信息><分隔符><后缀信息>
	// 错误信息输出格式
	ErrorBaseFormat   = "\x1b[31m%s\x1b[0m\n"     // 基础错误: <错误信息>
	ErrorSuffixFormat = "\x1b[31m%s\x1b[0m%s%s\n" // 带后缀的错误: <错误信息><分隔符><后缀信息>
)

// ---------- 环境变量

// 操作系统
var Platform = runtime.GOOS

// 用户名，当程序提权运行时，使用SUDO_USER变量获取提权前的用户名
var UserName = func() string {
	if GetVariable("SUDO_USER") != "" {
		return GetVariable("SUDO_USER")
	}
	return GetVariable("USER")
}()

// 用户信息
var UserInfo, _ = GetUserInfoByName(UserName)

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
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// SetVariable 设置环境变量
func SetVariable(key, value string) error {
	return os.Setenv(key, value)
}

// GetUserInfoByName 根据用户名获取用户信息
func GetUserInfoByName(userName string) (*user.User, error) {
	userInfo, err := user.Lookup(userName)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetUserInfoById 根据 ID 获取用户信息
func GetUserInfoById(userId int) (*user.User, error) {
	userInfo, err := user.LookupId(strconv.Itoa(userId))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// GetCurrentUserInfo 获取当前用户信息
func GetCurrentUserInfo() (*user.User, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}
