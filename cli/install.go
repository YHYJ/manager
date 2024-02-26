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

	"github.com/cheggaaa/pb/v3"
	"github.com/go-git/go-git/v5"
	"github.com/yhyj/manager/general"
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

// DownloadFile 通过 HTTP 协议下载文件
//
// 参数：
//   - url: 文件下载地址
//   - outputFile: 下载文件保存路径
//
// 返回：
//   - 错误信息
func DownloadFile(url string, outputFile string) error {
	// 发送GET请求并获取响应
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error sending download request: %s", err)
	}
	defer resp.Body.Close()

	// 检查返回值状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error downloading file: %s", resp.Status)
	}

	// 创建下载文件夹
	dir := filepath.Dir(outputFile)
	if dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("Error creating download folder: %s", err)
		}
	}

	// 创建本地文件
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("Error creating download file: %s", err)
	}
	defer file.Close()

	// 创建进度条模板
	barTemplate := `{{blue "Downloading:"}} {{bar . "[" "-" ">" " " "]"}} {{percent . "%.01f%%" "?"}} {{counters . "%s/%s" "%s/?" | green}} {{speed . | yellow}}`
	// 使用自定义模板创建进度条
	bar := pb.ProgressBarTemplate(barTemplate).Start64(resp.ContentLength)
	bar.Set(pb.Bytes, true)
	// 使用代理读取响应主体
	reader := bar.NewProxyReader(resp.Body)

	// 将响应主体复制到文件
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("Error writing download file: %s", err)
	}

	// 完成进度条
	bar.Finish()

	return nil
}

// InstallFile 安装文件，覆盖已存在的同名文件
//
// 参数：
//   - sourceFile: 源文件路径
//   - targetFile: 目标文件路径
//   - perm: 目标文件权限
//
// 返回：
//   - 错误信息
func InstallFile(sourceFile, targetFile string, perm os.FileMode) error {
	// 打开源文件
	sFile, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer sFile.Close()

	// 创建或打开目标文件
	tFile, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
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

// FileVerification 使用校验和文件校验文件的完整性
//
// 参数：
//   - checksumFile: 校验和文件
//   - filePath: 待校验文件
//
// 返回：
//   - 校验结果
//   - 错误信息
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
			continue
		}
		return false, fmt.Errorf("Checksum file format error, it should be: <checksum> <filename>")
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("Error reading checksum file: %s", err)
	}

	return false, nil
}
