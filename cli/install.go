/*
File: install.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-14 14:32:16

Description: 子命令`install`的实现
*/

package cli

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/yhyj/manager/general"
)

// CloneRepoViaHTTP 通过 HTTP 协议克隆仓库
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

// DownloadFile 通过 HTTP 协议下载文件
func DownloadFile(url string, outputFile string) (string, error) {
	// 发送GET请求并获取响应
	resp, err := http.Get(url)
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

// FileVerification 校验文件
func FileVerification(checksumFile, filePath string) (bool, error) {
	// 检查校验文件是否存在
	if !general.FileExist(checksumFile) {
		return false, fmt.Errorf("File %s does not exist", checksumFile)
	}
	// 检查待校验文件是否存在
	if !general.FileExist(filePath) {
		return false, fmt.Errorf("File %s does not exist", filePath)
	}

	// 打开校验文件
	file, err := os.Open(checksumFile)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// 扫描处理校验文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// 按行获取校验文件内容
		line := scanner.Text()
		// 以空格分割行
		lineParts := strings.Fields(line)

		if len(lineParts) == 2 {
			expectedChecksum := lineParts[0] // 期望的校验和
			filename := lineParts[1]         // 文件名

			// 检测校验文件中是否记载了指定文件的校验和信息
			if filename == filepath.Base(filePath) {
				// 计算文件的实际校验和
				actualChecksum, err := general.FileSHA256(filePath)
				if err != nil {
					return false, err
				}

				// 比对校验和
				if actualChecksum == expectedChecksum {
					return true, nil
				} else {
					return false, nil
				}
			}
			return false, fmt.Errorf("There is no %s info in the checksum file", filePath)
		}
		return false, fmt.Errorf("Checksum file format error, it should be: <checksum> <filename>")
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("Error reading checksum file: %s", err)
	}

	return false, nil
}
