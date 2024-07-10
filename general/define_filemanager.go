/*
File: define_filemanager.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:04:26

Description: 文件管理
*/

package general

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"sync"
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
			if _, err := file1.ReadAt(buffer1, offset); err != nil && err != io.EOF {
				errCh <- err
				return
			}

			// 从文件中读取指定大小的内容到 buffer
			if _, err = file2.ReadAt(buffer2, offset); err != nil && err != io.EOF {
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
