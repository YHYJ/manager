//go:build darwin

/*
File: setup.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 14:13:47

Description: 子命令 'setup' 的实现
*/

package cli

import "github.com/yhyj/manager/general"

// ProgramConfigurator 程序配置器
// 参数：
//   - flags: 系统信息各部分的开关
func ProgramConfigurator(flags map[string]bool) {
	// 配置 chezmoi
	if flags["chezmoiFlag"] {
		general.SetupChezmoi()
	}

	// 配置 cobra
	if flags["cobraFlag"] {
		general.SetupCobra()
	}

	// 配置 git
	if flags["gitFlag"] {
		general.SetupGit()
	}

	// 配置 golang
	if flags["goFlag"] {
		general.SetupGolang()
	}

	// 配置 pip
	if flags["pipFlag"] {
		general.SetupPip()
	}
}
