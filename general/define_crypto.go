/*
File: define_crypto.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-07-10 10:09:33

Description: 校验和加密
*/

package general

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileSHA256 计算文件的 SHA-256 校验和
//
// 参数：
//   - filePath: 待校验文件
//
// 返回：
//   - 校验和
//   - 错误信息
func FileSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// 计算校验和并将其转换为十六进制字符串
	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
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
	if !FileExist(checksumFile) {
		return false, fmt.Errorf("File %s does not exist", checksumFile)
	}
	// 检查待校验文件是否存在
	if !FileExist(filePath) {
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
				actualChecksum, err := FileSHA256(filePath)
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
