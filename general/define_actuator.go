/*
File: define_actuator.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 15:01:47

Description: 执行系统命令
*/

package general

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// RunCommand 运行命令，获取其标准输出和标准错误
//
//   - 标准输出和标准错误末尾自带的换行符已去除
//
// 参数：
//   - command: 命令
//   - args: 命令参数（每个以空格分隔的参数作为切片的一个元素）
//
// 返回：
//   - 标准输出
//   - 标准错误
//   - 错误信息
func RunCommand(command string, args []string) (string, string, error) {
	if _, err := exec.LookPath(command); err != nil {
		return "", "", err
	}

	// 定义命令
	cmd := exec.Command(command, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// 定义标准输入/输出/错误
	cmd.Stdin = os.Stdin
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 执行命令
	err := cmd.Run()

	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}
