/*
File: define_manager.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-16 15:26:20

Description: 安装/卸载管理
*/

package general

import (
	"io"
	"os"
)

// Install 安装，覆盖已存在的同名文件
//
// 参数：
//   - sourceFile: 源文件路径
//   - targetFile: 目标文件路径
//   - perm: 目标文件权限
//
// 返回：
//   - 错误信息
func Install(sourceFile, targetFile string, perm os.FileMode) error {
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

// Uninstall 卸载文件，自动检测文件是否存在
//
// 参数：
//   - targetFile: 目标文件路径
//
// 返回：
//   - 错误信息
func Uninstall(targetFile string) error {
	return DeleteFile(targetFile)
}

// InitPocketFile 初始化记账文件
//
// 参数：
//   - pocketFile: 记账文件路径
//
// 返回：
//   - 错误信息
func InitPocketFile(pocketFile string) error {
	if err := CreateFile(pocketFile); err != nil {
		return err
	}
	return EmptyFile(pocketFile)
}
