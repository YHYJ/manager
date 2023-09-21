/*
File: file_operation.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-08 16:04:26

Description: 文件操作
*/

package function

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// 判断文件是否存在
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断文件是否为空
func FileEmpty(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return true
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return true
	}
	return fi.Size() == 0
}

// 创建文件，如果其父目录不存在则创建父目录
func CreateFile(filePath string) error {
	if FileExist(filePath) {
		return nil
	}
	// 截取filePath的父目录
	parentPath := filePath[:strings.LastIndex(filePath, "/")]
	if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
		return err
	}
	_, err := os.Create(filePath)
	return err
}

// 创建文件夹，如果其父目录不存在则创建父目录
func CreateDir(dirPath string) error {
	if FileExist(dirPath) {
		return nil
	}
	// 截取dirPath的父目录
	parentPath := dirPath[:strings.LastIndex(dirPath, "/")]
	if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
		return err
	}
	return os.Mkdir(dirPath, os.ModePerm)
}

// 进到指定目录
func GoToDir(dirPath string) error {
	return os.Chdir(dirPath)
}

// 写入内容到文件
func WriteFile(filePath string, content string) error {
	// 文件存在
	if FileExist(filePath) {
		// 文件内容为空
		if FileEmpty(filePath) {
			// 打开文件并写入内容
			file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				return err
			} else {
				_, err := file.WriteString(content)
				if err != nil {
					return err
				}
			}
		} else {
			// 文件内容不为空
			return fmt.Errorf("File %s is not empty", filePath)
		}
	} else {
		// 文件不存在
		// 创建文件
		if err := CreateFile(filePath); err != nil {
			return err
		}
		// 打开文件并写入内容
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return err
		} else {
			_, err := file.WriteString(content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 删除文件
func DeleteFile(filePath string) error {
	if !FileExist(filePath) {
		return nil
	}
	return os.Remove(filePath)
}

// 并发比较两个文件是否相同
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

	// 如果文件大小不同，则直接返回不同
	if file1Size != file2Size {
		return false, nil
	}

	// 文件大小相同则（按块）比较文件内容
	const chunkSize = 1024 * 1024                             // 每次比较的块大小（1MB）
	numChunks := int((file1Size + chunkSize - 1) / chunkSize) // 计算文件需要被分成多少块

	equal := true                // 文件是否相同的标志位
	var wg sync.WaitGroup        // wg用于等待所有的goroutine执行完毕，然后关闭errCh通道
	errCh := make(chan error, 1) // errCh用于接收goroutine执行过程中返回的错误

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

			// 创建两个大小为size的buffer
			buffer1 := make([]byte, size)
			buffer2 := make([]byte, size)

			// 从文件中读取指定大小的内容到buffer
			_, err := file1.ReadAt(buffer1, offset)
			if err != nil && err != io.EOF {
				errCh <- err
				return
			}

			// 从文件中读取指定大小的内容到buffer
			_, err = file2.ReadAt(buffer2, offset)
			if err != nil && err != io.EOF {
				errCh <- err
				return
			}

			// 比较两个buffer是否相同
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

// 比较两个文件的内容是否相同
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
