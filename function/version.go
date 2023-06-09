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
var (
	name    string = "Manager"
	version string = "v0.0.6"
)

func ProgramInfo() string {
	programInfo := fmt.Sprintf("\033[1m%s\033[0m %s \033[1m%s\033[0m\n", name, "version", version)
	return programInfo
}
