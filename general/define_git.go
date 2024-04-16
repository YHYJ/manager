/*
File: define_git.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-16 15:21:32

Description: git 操作
*/

package general

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// CloneRepoViaHTTP 通过 HTTP 协议克隆仓库
//
// 参数：
//   - path: 本地仓库存储路径
//   - url: 远程仓库地址（不包括仓库名，https://github.com/{UserName}）
//   - repo: 仓库名
//
// 返回：
//   - 错误信息
func CloneRepoViaHTTP(path string, url string, repo string) error {
	_, err := git.PlainClone(filepath.Join(path, repo), false, &git.CloneOptions{
		URL:               url + "/" + repo,
		RecurseSubmodules: 1,
	})
	if err != nil {
		return err
	}
	return nil
}
