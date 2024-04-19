/*
File: define_other.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-11-24 13:35:18

Description: 处理一些杂事
*/

package general

import (
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
