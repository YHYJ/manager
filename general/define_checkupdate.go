/*
File: define_checkupdate.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-09-19 14:30:34

Description: 执行检查程序更新操作
*/

package general

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

// RequestApi 请求 API，返回响应体
func RequestApi(url string) ([]byte, error) {
	// 创建一个HTTP请求客户端
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	// 创建GET请求并设置请求头
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	// 发送HTTP请求并接收返回值
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 释放资源
	defer resp.Body.Close()

	// 检查返回值状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed with status: %s", resp.Status)
	}

	// 读取并解析响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	return body, nil
}

// GetLatestSourceTag 解析 API 响应体，获取源代码的最新Tag
//
//   - 该函数解析的是 https://api.github.com/repos/{OWNER}/{REPO}/tags 的返回值
//   - 用于通过Source安装程序时获取最新版本的Tag
func GetLatestSourceTag(body []byte) (string, error) {
	// 解码JSON格式的返回值
	var datas interface{}
	if err := json.Unmarshal(body, &datas); err != nil {
		return "", err
	}

	// 判断数据类型
	kind := reflect.ValueOf(datas).Kind()

	if kind == reflect.Slice { // '[{}]'结构
		// 判断响应体长度
		length := len(datas.([]interface{}))
		if length == 0 {
			return "", fmt.Errorf("Response body is empty")
		}
		// 获取最新版本信息，适用于版本信息和Tag信息同步的
		latestVersion := datas.([]interface{})[0].(map[string]interface{})["name"].(string)
		return latestVersion, nil
	} else {
		return "", fmt.Errorf("Response body has unknown structure")
	}
}

// GetLatestSourceHash 解析 API 响应体，获取源代码的最新提交的Hash
//
//   - 该函数解析的是 https://api.github.com/repos/{OWNER}/{REPO}/tags 的返回值
//   - 用于通过Source安装不带Tag的程序时获取最新版本的Hash
func GetLatestSourceHash(body []byte) (string, error) {
	// 解码JSON格式的返回值
	var datas interface{}
	if err := json.Unmarshal(body, &datas); err != nil {
		return "", err
	}

	// 判断数据类型
	kind := reflect.ValueOf(datas).Kind()

	if kind == reflect.Map { // '{}'结构
		// 判断响应体长度
		length := len(datas.(map[string]interface{}))
		if length == 0 {
			return "", fmt.Errorf("Response body is empty")
		}
		// 获取文件哈希值，适用于不带外部版本信息的
		fileHash := datas.(map[string]interface{})["sha"].(string)
		return fileHash, nil
	} else {
		return "", fmt.Errorf("Response body has unknown structure")
	}
}

// GetLatestReleaseTag 解析 API 响应体，获取Release的最新Tag
//
//   - 该函数解析的是 https://api.github.com/repos/{OWNER}/{REPO}/releases/latest 的返回值
//   - 用于通过Release安装程序时获取最新版本的Tag
func GetLatestReleaseTag(body []byte) (string, error) {
	// 解码JSON格式的返回值
	var datas interface{}
	if err := json.Unmarshal(body, &datas); err != nil {
		return "", err
	}

	// 判断数据类型
	kind := reflect.ValueOf(datas).Kind()

	if kind == reflect.Map { // '{}'结构
		// 判断响应体长度
		length := len(datas.(map[string]interface{}))
		if length == 0 {
			return "", fmt.Errorf("Response body is empty")
		}
		// 获取最新Release对应的Tag
		latestTag := datas.(map[string]interface{})["tag_name"].(string)
		return latestTag, nil
	} else {
		return "", fmt.Errorf("Response body has unknown structure")
	}
}

// 指定要获取信息的文件名
type FileName struct {
	ChecksumsFile string `json:"checksums"`
	ArchiveFile   string `json:"archive"`
}

// 存储多文件信息
type multipleFilesInfo struct {
	ChecksumsFileInfo singleFileInfo `json:"checksums_file_info"`
	ArchiveFileInfo   singleFileInfo `json:"archive_file_info"`
}

// 存储单文件信息
type singleFileInfo struct {
	Name          string  `json:"name"`
	Size          float64 `json:"size"`
	ContentType   string  `json:"content_type"`
	DownloadUrl   string  `json:"download_url"`
	DownloadCount float64 `json:"download_count"`
}

// GetReleaseFileInfo 解析 API 响应体，获取Release文件的信息
//
//   - 该函数解析的是 https://api.github.com/repos/{OWNER}/{REPO}/releases/latest 的返回值
//   - 用于通过Release安装程序时获取校验文件、压缩包等文件的信息
func GetReleaseFileInfo(body []byte, fileName FileName) (multipleFilesInfo, error) {
	filesInfo := multipleFilesInfo{}      // 存储多文件信息
	checksumsFileInfo := singleFileInfo{} // 存储校验文件信息
	archiveFileInfo := singleFileInfo{}   // 存储压缩包信息

	// 解码JSON格式的返回值
	var datas interface{}
	if err := json.Unmarshal(body, &datas); err != nil {
		return filesInfo, err
	}

	// 判断数据类型
	kind := reflect.ValueOf(datas).Kind()

	if kind == reflect.Map { // '{}'结构
		// 判断响应体长度
		length := len(datas.(map[string]interface{}))
		if length == 0 {
			return filesInfo, fmt.Errorf("Response body is empty")
		}
		// 获取最新Release的Assets信息，下载链接等包含在里面
		assets := datas.(map[string]interface{})["assets"].([]interface{})
		for _, asset := range assets {
			if asset.(map[string]interface{})["name"] == fileName.ChecksumsFile {
				checksumsFileInfo.Name = asset.(map[string]interface{})["name"].(string)
				checksumsFileInfo.Size = asset.(map[string]interface{})["size"].(float64)
				checksumsFileInfo.ContentType = asset.(map[string]interface{})["content_type"].(string)
				checksumsFileInfo.DownloadUrl = asset.(map[string]interface{})["browser_download_url"].(string)
				checksumsFileInfo.DownloadCount = asset.(map[string]interface{})["download_count"].(float64)
			}
			if asset.(map[string]interface{})["name"] == fileName.ArchiveFile {
				archiveFileInfo.Name = asset.(map[string]interface{})["name"].(string)
				archiveFileInfo.Size = asset.(map[string]interface{})["size"].(float64)
				archiveFileInfo.ContentType = asset.(map[string]interface{})["content_type"].(string)
				archiveFileInfo.DownloadUrl = asset.(map[string]interface{})["browser_download_url"].(string)
				archiveFileInfo.DownloadCount = asset.(map[string]interface{})["download_count"].(float64)
			}
		}
		filesInfo.ChecksumsFileInfo = checksumsFileInfo
		filesInfo.ArchiveFileInfo = archiveFileInfo
		return filesInfo, nil
	} else {
		return filesInfo, fmt.Errorf("Response body has unknown structure")
	}
}
