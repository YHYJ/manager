/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-14 14:32:16

Description: 子命令`install`的实现
*/

package function

import (
	"fmt"
	"io"
	"os"

	"github.com/go-git/go-git/v5"
)

// 通过HTTP协议克隆仓库
func CloneRepoViaHTTP(path string, url string, repo string) {
	_, err := git.PlainClone(path+"/"+repo, false, &git.CloneOptions{
		URL:               url + "/" + repo,
		RecurseSubmodules: 1,
	})
	if err != nil {
		fmt.Printf("\x1b[31m==>\x1b[0m Clone \x1b[32m%s\x1b[0m error: %s\n", repo, err)
	} else {
		fmt.Printf("\x1b[32;1m==>\x1b[0m Clone \x1b[32m%s\x1b[0m success\n", repo)
	}
}

// 复制文件，如果文件存在则覆盖
func InstallFile(sourceFile, targetFile string) error {
	// 打开源文件
	sFile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer sFile.Close()

	// 创建或打开目标文件
	tFile, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer tFile.Close()

	// 复制文件内容
	_, err = io.Copy(tFile, sFile)
	if err != nil {
		return err
	}
	return nil
}
