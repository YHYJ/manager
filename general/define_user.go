/*
File: define_user.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-06-12 16:57:10

Description: 获取用户信息
*/

package general

import "os/user"

// GetUserName 获取用户真实/显示名称，可能为空
//
// 返回：
//   - 用户名称
//   - 错误信息
func GetUserName() (string, error) {
	userData, err := user.Current() // 获取用户信息
	if err != nil {
		return "", err
	}
	return userData.Name, nil
}
