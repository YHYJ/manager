/*
File: define_other.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-11-24 13:35:18

Description: 处理一些杂事
*/

package general

import (
	"bufio"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gookit/color"
)

// RealLength 去除控制字符和图标的附加字符，获取文本实际长度
//
// 参数：
//   - text: 文本
//
// 返回：
//   - 实际长度
func RealLength(text string) int {
	combinedRegex := regexp.MustCompile(`\x1b\[[0-9;]*m|\x{FE0E}|\x{FE0F}`)
	return utf8.RuneCountInString(combinedRegex.ReplaceAllString(text, ""))
}

// PrintDelimiter 打印分隔符
//
// 参数：
//   - length: 分隔符长度
func PrintDelimiter(length int) {
	dashes := strings.Repeat("-", length)     // 组装分隔符
	color.Printf("%s\n", FgBlackText(dashes)) // 美化输出
}

// Delay 延时
//
// 参数：
//   - second: 延时秒数
func Delay(second float32) {
	time.Sleep(time.Duration(second*1000) * time.Millisecond)
}

// AskUser 询问用户
//
// 参数：
//   - question: 问题
//   - answer: 期望的回答（每个选项之间用斜线/隔开），例如 "y/N" 代表期望输入 y 或 n，其中大写字母代表默认值，示例是 N
//
// 返回：
//   - 用户的回答
//   - 错误信息
func AskUser(question string, answer string) (string, error) {
	color.Printf("%s [%s] ", question, answer)

	// 从标准输入中读取用户的回答
	reader := bufio.NewReader(os.Stdin)
	userAnswer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// 将用户的回答转换为小写并去除首尾空格
	userAnswer = strings.TrimSpace(strings.ToLower(userAnswer))

	return userAnswer, nil
}

// GetInput 获取用户输入
//
// 参数：
//   - tips: 提示信息
//   - default: 用户未输入时的默认值
//
// 返回：
//   - 用户输入
//   - 错误信息
func GetInput(tips string, defaultValue string) (string, error) {
	color.Printf("%s %s: ", tips, SecondaryText(color.Sprintf("(%s)", defaultValue)))

	// 从标准输入中读取用户的回答
	reader := bufio.NewReader(os.Stdin)
	originalValue, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	value := func() string {
		if len(originalValue) <= 1 {
			return defaultValue
		}
		return originalValue
	}()

	return value, nil
}
