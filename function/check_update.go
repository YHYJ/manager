/*
File: check_update.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-09-19 14:30:34

Description: 执行检查程序更新操作
*/

package function

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Data struct {
	Name string `json:"name"`
}

// 获取最新版本信息，适用于版本信息和Tag信息同步的
func GetLatestVersion(url string) (error, string) {
	// 创建一个HTTP请求客户端
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	// 创建GET请求并设置请求头
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err, ""
	}
	req.Header.Set("Accept", "application/json")
	// 发送HTTP请求并接收返回值
	resp, err := client.Do(req)
	if err != nil {
		return err, ""
	}

	// 释放资源
	defer resp.Body.Close()

	// 检查返回值状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failed with status: %s", resp.Status), ""
	}
	// 解码JSON格式的返回值
	var datas []Data
	err = json.NewDecoder(resp.Body).Decode(&datas)
	if err != nil {
		return err, ""
	}
	// 获取其中的Version信息
	latestVersion := datas[0].Name

	return nil, latestVersion
}
