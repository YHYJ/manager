/*
File: variable_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:01:45

Description: 操作变量
*/

package function

import (
	"os"
	"os/user"
	"runtime"
	"strconv"
)

var platformChart = map[string]map[string]string{
	"linux": {
		"HOME": "HOME",
		"PWD":  "PWD",
		"USER": "USER",
	},
	"darwin": {
		"HOME": "HOME",
		"PWD":  "PWD",
		"USER": "USER",
	},
	"windows": {
		"HOME": "USERPROFILE",
		"PWD":  "PWD",
		"USER": "USERNAME",
	},
}

var platform = runtime.GOOS

// 获取环境变量
func GetVariable(key string) string {
	varKey := platformChart[platform][key]
	return os.Getenv(varKey)
}

// 获取不在环境变量中的HOSTNAME
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}

// 设置环境变量
func SetVariable(key, value string) error {
	return os.Setenv(key, value)
}

// 根据ID获取用户信息
func GetUserInfo(uid int) (*user.User, error) {
	userInfo, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

// 获取当前用户信息
func GetCurrentUserInfo() (*user.User, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	return currentUser, nil
}
