/*
File: file_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:04:26

Description: 文件操作
*/

package function

import (
	"fmt"
	"os"
	"strings"
)

// 判断文件是否存在
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断文件是否为空
func FileEmpty(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return true
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return true
	}
	return fi.Size() == 0
}

// 创建文件，如果其父目录不存在则创建父目录
func CreateFile(filePath string) error {
	if FileExist(filePath) {
		return nil
	}
	// 截取filePath的父目录
	dirPath := filePath[:strings.LastIndex(filePath, "/")]
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}
	_, err := os.Create(filePath)
	return err
}

// 写入内容到文件
func WriteFile(filePath string, content string) error {
	// 文件存在
	if FileExist(filePath) {
		// 文件内容为空
		if FileEmpty(filePath) {
			// 打开文件并写入内容
			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				return err
			} else {
				_, err := file.WriteString(content)
				if err != nil {
					return err
				}
			}
		} else {
			// 文件内容不为空
			return fmt.Errorf("file %s is not empty", filePath)
		}
	} else {
		// 文件不存在
		// 创建文件
		if err := CreateFile(filePath); err != nil {
			return err
		}
		// 打开文件并写入内容
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return err
		} else {
			_, err := file.WriteString(content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 删除文件
func DeleteFile(filePath string) error {
	if !FileExist(filePath) {
		return nil
	}
	return os.Remove(filePath)
}
