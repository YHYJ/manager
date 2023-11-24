/*
File: define_other.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-11-24 13:35:18

Description: 处理一些杂事
*/

package general

import (
	"fmt"
	"regexp"
	"strings"
)

// RealLength 去除转义字符，获取文本实际长度
//
// 参数：
//   - text: 文本
//
// 返回：
//   - 实际长度
func RealLength(text string) int {
	controlRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return len(controlRegex.ReplaceAllString(text, ""))
}

func PrintDelimiter(length int) {
	dashes := strings.Repeat("-", length-1) // 组装分隔符（减去行尾换行符的一个长度）
	fmt.Printf(LineHiddenFormat, dashes)    // 美化输出
}
