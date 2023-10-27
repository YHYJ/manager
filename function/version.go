/*
File: version.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-07 11:08:32

Description: 子命令`version`功能函数
*/

package function

import "fmt"

// 程序信息
const (
	Name    = "Manager"
	Version = "v0.8.7"
	Path    = "github.com/yhyj/manager"
)

func ProgramInfo(only bool) string {
	programInfo := fmt.Sprintf("%s\n", Version)
	if !only {
		programInfo = fmt.Sprintf("%s version %s\n", Name, Version)
	}
	return programInfo
}
