/*
File: config.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:05:47

Description: 子命令`config`的实现
*/

package function

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

// 判断文件是不是toml文件
func isTomlFile(filePath string) bool {
	if strings.HasSuffix(filePath, ".toml") {
		return true
	}
	return false
}

// 读取toml配置文件
func GetTomlConfig(filePath string) (*toml.Tree, error) {
	if !FileExist(filePath) {
		return nil, fmt.Errorf("open %s: no such file or directory", filePath)
	}
	if !isTomlFile(filePath) {
		return nil, fmt.Errorf("open %s: is not a toml file", filePath)
	}
	tree, err := toml.LoadFile(filePath)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// 写入toml配置文件
func WriteTomlConfig(filePath string) (int64, error) {
	// 获取指定用户信息
	userInfo, err := GetUserInfo(1000)
	if err != nil {
		return 0, err
	}
	// 定义一个map[string]interface{}类型的变量并赋值
	exampleConf := map[string]interface{}{
		"variable": map[string]interface{}{
			"http_proxy": "http://127.0.0.1:1080",
			"https_proxy": "http://127.0.0.1:1080",
		},
		"install": map[string]interface{}{
			"path": "/usr/local/bin",
			"temp": "/tmp/manager-build",
			"go": map[string]interface{}{
				"source": "https://github.com/YHYJ",
				"fallback_source": "https://git.yj1516.top/YJ",
				"names":  []string{"checker", "clone-repos", "eniac", "kbdstage", "manager", "rolling", "scleaner", "skynet"},
				"completion_dir": userInfo.HomeDir + "/.cache/oh-my-zsh/completions",
			},
			"shell": map[string]interface{}{
				"source": "https://github.com/YHYJ",
				"fallback_source": "https://git.yj1516.top/YJ",
				"repo": "Program",
				"dir":    "/System-Script/app",
				"names":  []string{"collect-system", "configure-dtags", "py-virtualenv-tool", "save-docker-images", "sfm", "spacevim-update", "spider", "system-checkupdates", "trash-manager", "usb-manager"},
			},
		},
	}
	if !FileExist(filePath) {
		return 0, fmt.Errorf("open %s: no such file or directory", filePath)
	}
	if !isTomlFile(filePath) {
		return 0, fmt.Errorf("open %s: is not a toml file", filePath)
	}
	// 把exampleConf转换为*toml.Tree
	tree, err := toml.TreeFromMap(exampleConf)
	if err != nil {
		return 0, err
	}
	// 打开一个文件并获取io.Writer接口
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	return tree.WriteTo(file)
}
