/*
File: define_download.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-04-16 15:23:30

Description: 文件下载
*/

package general

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

// DownloadFile 通过 HTTP 协议下载文件
//
// 参数：
//   - url: 文件下载地址
//   - outputFile: 下载文件保存路径
//   - progressParameters: 进度条参数
//
// 返回：
//   - 错误信息
func DownloadFile(url string, outputFile string, progressParameters map[string]string) error {
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

	if progressParameters["view"] == "0" {
		// 将响应主体复制到文件
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("Error writing download file: %s", err)
		}
	} else {
		// 创建进度条模板
		barTemplate := `{{string . "action" | green}} {{string . "prefix"}} {{string . "project" | blue}} {{string . "sep" | blue}} {{string . "fileName" | blue}} {{string . "suffix"}} {{bar . "[" "-" ">" " " "]"}} {{percent . "%.01f%%" "?"}} {{counters . "%s/%s" "%s/?" | green}} {{speed . | yellow}}`
		// 使用自定义模板创建进度条
		bar := pb.ProgressBarTemplate(barTemplate).Start64(resp.ContentLength)
		bar.Set(pb.Bytes, true)
		bar.Set("action", progressParameters["action"]).Set("prefix", progressParameters["prefix"]).Set("project", progressParameters["project"]).Set("sep", progressParameters["sep"]).Set("fileName", progressParameters["fileName"]).Set("suffix", progressParameters["suffix"])
		// 使用代理读取响应主体
		reader := bar.NewProxyReader(resp.Body)

		// 将响应主体复制到文件
		_, err = io.Copy(file, reader)
		if err != nil {
			return fmt.Errorf("Error writing download file: %s", err)
		}

		// 完成进度条
		bar.Finish()
	}

	return nil
}
