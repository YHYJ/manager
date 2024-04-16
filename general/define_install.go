/*
File: define_install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-16 15:26:20

Description: 文件安装
*/

package general

import (
	"io"
	"os"
)

// InstallFile 安装文件，覆盖已存在的同名文件
//
// 参数：
//   - sourceFile: 源文件路径
//   - targetFile: 目标文件路径
//   - perm: 目标文件权限
//
// 返回：
//   - 错误信息
func InstallFile(sourceFile, targetFile string, perm os.FileMode) error {
	// 打开源文件
	sFile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer sFile.Close()

	// 创建或打开目标文件
	tFile, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer tFile.Close()

	// 复制文件内容
	_, err = io.Copy(tFile, sFile)
	if err != nil {
		return err
	}
	return nil
}
