/*
File: define_window.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-06-19 10:16:34

Description: 窗口操作
*/

package general

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// GetTerminalSize 获取终端窗口大小
//
// 返回：
//   - 窗口宽度
//   - 窗口高度
//   - 错误信息
func GetTerminalSize() (int, int, error) {
	// 获取标准输入的文件描述符
	fd := int(os.Stdout.Fd())

	// 检查文件描述符是否指向终端
	if !term.IsTerminal(fd) {
		return 0, 0, fmt.Errorf("File descriptor %d is not a terminal", fd)
	}

	// 获取终端窗口大小
	width, height, err := term.GetSize(fd)
	if err != nil {
		return 0, 0, err
	}

	return width, height, nil
}
