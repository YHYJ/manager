/*
File: define_archive.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-07-10 10:01:45

Description: 处理归档文件
*/

package general

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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
		if err := extractZipFile(file, outputTo); err != nil {
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

		if err = extractTarFile(header, tarReader, outputTo); err != nil {
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
