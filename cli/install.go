/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-14 14:32:16

Description: 子命令`install`的实现
*/

package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
)

// CloneRepoViaHTTP 通过 HTTP 协议克隆仓库
func CloneRepoViaHTTP(path string, url string, repo string) error {
	_, err := git.PlainClone(path+"/"+repo, false, &git.CloneOptions{
		URL:               url + "/" + repo,
		RecurseSubmodules: 1,
	})
	if err != nil {
		return err
	}
	return nil
}

// DownloadFile 通过 HTTP 协议下载文件
func DownloadFile(url string, urlFile string, outputFile string) (string, error) {
	fileDownloadUrl := fmt.Sprintf("%s/%s", url, urlFile)

	// 发送GET请求并获取响应
	resp, err := http.Get(fileDownloadUrl)
	if err != nil {
		return "", fmt.Errorf("Error sending download request: %s", err)
	}
	defer resp.Body.Close()

	// 检查返回值状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error downloading file: %s", resp.Status)
	}

	// 创建下载文件夹
	dir, filename := filepath.Split(outputFile)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("Error creating download folder: %s", err)
		}
	}

	// 创建本地文件
	outputFileFullname := filepath.Join(dir, filename)
	file, err := os.Create(outputFileFullname)
	if err != nil {
		return "", fmt.Errorf("Error creating download file: %s", err)
	}
	defer file.Close()

	// 将响应主体复制到文件
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error writing download file: %s", err)
	}

	return outputFileFullname, nil
}

// InstallFile 安装文件，如果文件存在则覆盖
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
