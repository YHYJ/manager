/*
File: define_filemanager.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:04:26

Description: 文件管理
*/

package general

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gookit/color"
)

// ReadFile 依次读取文件每行内容
//
// 参数：
//   - file: 文件路径
//
// 返回：
//   - 指定行的内容
func ReadFile(file string) ([]string, error) {
	// 打开文件
	text, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer text.Close()

	// 创建一个 Scanner 对象
	scanner := bufio.NewScanner(text)

	// 存储读取到的每行内容的切片
	var lines []string

	// 逐行读取文件内容
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 检查是否出现了读取错误
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// ReadFileLine 读取文件指定行
//
// 参数：
//   - file: 文件路径
//   - line: 行号
//
// 返回：
//   - 指定行的内容
func ReadFileLine(file string, line int) string {
	// 打开文件
	text, err := os.Open(file)
	if err != nil {
		fileName, lineNo := GetCallerInfo()
		color.Danger.Printf("Open file error (%s:%d): %s\n", fileName, lineNo+1, err)
	}
	defer text.Close()

	// 创建一个扫描器对象按行遍历
	scanner := bufio.NewScanner(text)
	// 行计数
	count := 1
	// 逐行读取，输出指定行
	for scanner.Scan() {
		if line == count {
			return scanner.Text()
		}
		count++
	}
	return ""
}

// ReadFileKey 读取文件包含关键字的行
//
// 参数：
//   - file: 文件路径
//   - key: 关键字
//
// 返回：
//   - 包含关键字的行的内容
func ReadFileKey(file, key string) string {
	// 打开文件
	text, err := os.Open(file)
	if err != nil {
		fileName, lineNo := GetCallerInfo()
		color.Danger.Printf("Open file error (%s:%d): %s\n", fileName, lineNo+1, err)
	}
	defer text.Close()

	// 创建一个扫描器对象按行遍历
	scanner := bufio.NewScanner(text)
	// 逐行读取，输出指定行
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), key) {
			return scanner.Text()
		}
	}
	return ""
}

// ReadFileCount 获取文件包含关键字的行的计数
//
// 参数：
//   - file: 文件路径
//   - key: 关键字
//
// 返回：
//   - 包含关键字的行的数量
func ReadFileCount(file, key string) int {
	// 打开文件
	text, err := os.Open(file)
	if err != nil {
		fileName, lineNo := GetCallerInfo()
		color.Danger.Printf("Open file error (%s:%d): %s\n", fileName, lineNo+1, err)
	}
	defer text.Close()

	// 创建一个扫描器对象按行遍历
	scanner := bufio.NewScanner(text)
	// 计数器
	count := 0
	// 逐行读取，输出指定行
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), key) {
			count++
		}
	}
	return count
}

// FileExist 判断文件是否存在
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 文件存在返回 true，否则返回 false
func FileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetAbsPath 获取指定文件的绝对路径
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 文件的绝对路径
func GetAbsPath(filePath string) string {
	// 获取绝对路径
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return ""
	} else {
		return absPath
	}
}

// GetFileDir 获取文件所在目录
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 文件所在目录
func GetFileDir(filePath string) string {
	return filepath.Dir(filePath)
}

// FileEmpty 判断文件是否为空
//
//   - 无法判断文件夹
//
// 参数：
//   - file: 文件路径
//
// 返回：
//   - 文件为空返回 true，否则返回 false
func FileEmpty(file string) bool {
	text, err := os.Open(file)
	if err != nil {
		return true
	}
	defer text.Close()

	fi, err := text.Stat()
	if err != nil {
		return true
	}
	return fi.Size() == 0
}

// EmptyFile 清空文件内容，文件不存在则创建
//
// 参数：
//   - file: 文件路径
//
// 返回：
//   - 错误信息
func EmptyFile(file string) error {
	// 打开文件，如果不存在则创建，文件权限为读写
	text, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer text.Close()

	// 清空文件内容
	if err := text.Truncate(0); err != nil {
		return err
	}
	return nil
}

// FolderEmpty 判断文件夹是否为空
//
//   - 包括隐藏文件
//
// 参数：
//   - dir: 文件夹路径
//
// 返回：
//   - 文件夹为空返回 true，否则返回 false
func FolderEmpty(dir string) bool {
	text, err := os.Open(dir)
	if err != nil {
		return true
	}
	defer text.Close()

	_, err = text.Readdir(1)
	if err == io.EOF {
		return true
	}
	return false
}

// ListFolderFiles 列出指定文件夹下的所有文件
//
// 参数：
//   - dir: 文件夹路径
//
// 返回：
//   - 文件列表
//   - 错误信息
func ListFolderFiles(dir string) ([]string, error) {
	files := []string{}

	// 打开文件夹
	text, err := os.Open(dir)
	if err != nil {
		return files, err
	}
	defer text.Close()

	// 读取文件夹中的文件
	fileInfos, err := text.ReadDir(-1)
	if err != nil {
		return files, err
	}

	// 遍历文件夹中的文件
	for _, fileInfo := range fileInfos {
		// 判断是否为文件
		if !fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}

	return files, nil
}

// CreateFile 创建文件，包括其父目录
//
// 参数：
//   - file: 文件路径
//
// 返回：
//   - 错误信息
func CreateFile(file string) error {
	if FileExist(file) {
		return nil
	}
	// 创建父目录
	parentPath := filepath.Dir(file)
	if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
		return err
	}
	// 创建文件
	if _, err := os.Create(file); err != nil {
		return err
	}

	return nil
}

// CreateDir 创建文件夹
//
// 参数：
//   - dir: 文件夹路径
//
// 返回：
//   - 错误信息
func CreateDir(dir string) error {
	if FileExist(dir) {
		return nil
	}
	return os.MkdirAll(dir, os.ModePerm)
}

// GoToDir 进到指定文件夹
//
// 参数：
//   - dirPath: 文件夹路径
//
// 返回：
//   - 错误信息
func GoToDir(dirPath string) error {
	return os.Chdir(dirPath)
}

// WriteFile 写入内容到文件，文件不存在则创建，不自动换行
//
// 参数：
//   - filePath: 文件路径
//   - content: 内容
//   - mode: 写入模式，追加('a', O_APPEND, 默认)或覆盖('t', O_TRUNC)
//
// 返回：
//   - 错误信息
func WriteFile(filePath, content, mode string) error {
	// 确定写入模式
	writeMode := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if mode == "t" {
		writeMode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}

	// 将内容写入文件
	file, err := os.OpenFile(filePath, writeMode, 0666)
	if err != nil {
		return err
	}
	if _, err = file.WriteString(content); err != nil {
		return err
	}
	return nil
}

// WriteFileWithNewLine 写入内容到文件，文件不存在则创建，自动换行
//
// 参数：
//   - filePath: 文件路径
//   - content: 写入内容
//   - mode: 写入模式，追加('a', O_APPEND, 默认)或覆盖('t', O_TRUNC)
//
// 返回：
//   - 错误信息
func WriteFileWithNewLine(filePath, content, mode string) error {
	// 确定写入模式
	writeMode := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	if mode == "t" {
		writeMode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}

	// 将内容写入文件
	file, err := os.OpenFile(filePath, writeMode, 0666)
	if err != nil {
		return err
	}
	if _, err = file.WriteString(content + "\n"); err != nil {
		return err
	}
	return nil
}

// DeleteFile 删除文件，如果目标是文件夹则包括其下所有文件
//
// 参数：
//   - filePath: 文件路径
//
// 返回：
//   - 错误信息
func DeleteFile(filePath string) error {
	if !FileExist(filePath) {
		return nil
	}
	return os.RemoveAll(filePath)
}

// CompareFile 并发比较两个文件是否相同
//
// 参数：
//   - file1Path: 文件1路径
//   - file2Path: 文件2路径
//
// 返回：
//   - 文件相同返回 true，出错或不同返回 false
func CompareFile(file1Path string, file2Path string) (bool, error) {
	// 尝试打开文件
	file1, err := os.Open(file1Path)
	if err != nil {
		return false, err
	}
	defer file1.Close()
	file2, err := os.Open(file2Path)
	if err != nil {
		return false, err
	}
	defer file2.Close()

	// 获取文件大小
	file1Info, err := file1.Stat()
	if err != nil {
		return false, err
	}
	file2Info, err := file2.Stat()
	if err != nil {
		return false, err
	}
	file1Size := file1Info.Size()
	file2Size := file2Info.Size()

	// 如果文件大小不同则直接返回
	if file1Size != file2Size {
		return false, nil
	}

	// 文件大小相同则（按块）比较文件内容
	const chunkSize = 1024 * 1024                             // 每次比较的块大小（1MB）
	numChunks := int((file1Size + chunkSize - 1) / chunkSize) // 计算文件需要被分成多少块

	equal := true                // 文件是否相同的标志位
	var wg sync.WaitGroup        // wg 用于等待所有的 goroutine 执行完毕，然后关闭 errCh 通道
	errCh := make(chan error, 1) // errCh 用于接收 goroutine 执行过程中返回的错误

	for i := 0; i < numChunks; i++ { // 同时比较多个块
		wg.Add(1)
		go func(chunkIndex int) {
			defer wg.Done()

			// 计算当前块的偏移量和大小
			offset := int64(chunkIndex) * chunkSize
			size := chunkSize
			if offset+int64(size) > file1Size {
				size = int(file1Size - offset)
			}

			// 创建两个大小为 size 的 buffer
			buffer1 := make([]byte, size)
			buffer2 := make([]byte, size)

			// 从文件中读取指定大小的内容到 buffer
			_, err := file1.ReadAt(buffer1, offset)
			if err != nil && err != io.EOF {
				errCh <- err
				return
			}

			// 从文件中读取指定大小的内容到 buffer
			_, err = file2.ReadAt(buffer2, offset)
			if err != nil && err != io.EOF {
				errCh <- err
				return
			}

			// 比较两个 buffer 是否相同
			if !bytesEqual(buffer1, buffer2) {
				equal = false
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return false, err
		}
	}

	return equal, nil
}

// bytesEqual 比较两个文件的内容
//
// 参数：
//   - b1: 文件1内容
//   - b2: 文件2内容
//
// 返回：
//   - 相同返回 true，不同返回 false
func bytesEqual(b1 []byte, b2 []byte) bool {
	if len(b1) != len(b2) {
		return false
	}

	for i := 0; i < len(b1); i++ {
		if b1[i] != b2[i] {
			return false
		}
	}

	return true
}

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

// UnzipFile 检测压缩文件类型，执行相应的解压函数
//
// 参数：
//   - filePath: 待解压文件
//   - outputDir: 解压文件的存储目录
//
// 返回：
//   - 错误信息
func UnzipFile(filePath, outputDir string) error {
	// 读取压缩文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 接收解压缩函数并执行
	var unzipFunc func(string, string) error

	// 通过读取文件头部信息和魔数对比来检测文件类型
	bufferedReader := bufio.NewReader(file)
	fileType, err := bufferedReader.Peek(10)
	if err != nil && err != io.EOF {
		return err
	}
	// 根据文件类型选择相应的解压缩函数
	switch {
	case strings.HasPrefix(string(fileType), "PK\x03\x04"): // zip 文件的魔数
		unzipFunc = unzipZip
	case strings.HasPrefix(string(fileType), "\x1F\x8B"): // tar.gz 文件的魔数
		unzipFunc = unzipTarGz
	default:
		return fmt.Errorf("Unsupported compressed file type")
	}

	return unzipFunc(filePath, outputDir)
}

// unzipZip 解压 zip 格式压缩包
//
// 参数：
//   - filePath: 待解压文件
//   - outputDir: 解压文件的存储目录
//
// 返回：
//   - 错误信息
func unzipZip(filePath, outputDir string) error {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	fileNameWithExt := filepath.Base(filePath)                         // 带后缀的文件名
	fileExt := filepath.Ext(fileNameWithExt)                           // 文件后缀
	fileNameWithoutExt := strings.TrimSuffix(fileNameWithExt, fileExt) // 不带后缀的文件名
	outputTo := filepath.Join(outputDir, fileNameWithoutExt)
	for _, file := range reader.File {
		err := extractZipFile(file, outputTo)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractZipFile 解压 zip 格式压缩包中的单个文件
//
// 参数：
//   - file: 待解压文件
//   - outputDir: 解压文件的存储目录
//
// 返回：
//   - 错误信息
func extractZipFile(file *zip.File, outputDir string) error {
	// 创建目标目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// 拼接 file 解压后的路径
	path := filepath.Join(outputDir, file.Name)

	// 如果 file 是文件夹
	if file.FileInfo().IsDir() {
		// 在输出目录创建该文件夹后返回
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
		return nil
	}

	// 如果 file 是普通文件
	// 1. 读取该文件内容
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	// 2. 在输出目录创建该文件
	writer, err := os.Create(path)
	if err != nil {
		return err
	}
	defer writer.Close()

	// 3. 拷贝 file 的内容到创建的文件
	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}

	// 4. 从 file 恢复文件权限
	if err := os.Chmod(path, file.Mode()); err != nil {
		return err
	}

	return nil
}

// unzipTarGz 解压 tar.gz 格式压缩包
//
// 参数：
//   - filePath: 待解压文件
//   - outputDir: 解压文件的存储目录
//
// 返回：
//   - 错误信息
func unzipTarGz(filePath, outputDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 gzip 读取文件
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// 使用 tar 读取 gzip 流
	tarReader := tar.NewReader(gzipReader)

	fileName := func() string {
		switch {
		case strings.HasSuffix(filePath, ".tar.gz"):
			return strings.TrimSuffix(filepath.Base(filePath), ".tar.gz")
		default:
			fileNameWithExt := filepath.Base(filePath)                         // 带后缀的文件名
			fileExt := filepath.Ext(fileNameWithExt)                           // 文件后缀
			fileNameWithoutExt := strings.TrimSuffix(fileNameWithExt, fileExt) // 不带后缀的文件名
			return fileNameWithoutExt
		}
	}()
	outputTo := filepath.Join(outputDir, fileName)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = extractTarFile(header, tarReader, outputTo)
		if err != nil {
			return err
		}
	}

	return nil
}

// extractTarFile 解压 tar 格式压缩包中的单个文件
//
// 参数：
//   - header: 待解压文件头信息
//   - reader: 待解压文件内容读取器
//   - outputDir: 解压文件的存储目录
//
// 返回：
//   - 错误信息
func extractTarFile(header *tar.Header, reader io.Reader, outputDir string) error {
	// 创建目标目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// 拼接 header.Name 解压后的路径
	path := filepath.Join(outputDir, header.Name)

	switch header.Typeflag {
	case tar.TypeDir: // 如果 header.Name 是文件夹
		// 在输出目录创建该文件夹
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	case tar.TypeReg: // 如果 header.Name 是普通文件
		// 因为 tar 包中 header.Name 是最终文件（即文件或空文件夹），所以需要在输出目录创建其父文件夹
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// 1. 在输出目录创建该文件
		writer, err := os.Create(path)
		if err != nil {
			return err
		}
		defer writer.Close()

		// 2. 拷贝 header.Name 的内容到创建的文件
		if _, err := io.Copy(writer, reader); err != nil {
			return err
		}

		// 3. 从 header.Name 恢复文件权限
		if err = os.Chmod(path, os.FileMode(header.Mode)); err != nil {
			return err
		}
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
