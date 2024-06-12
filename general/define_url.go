/*
File: define_url.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-06-12 14:31:51

Description: 处理 URL
*/

package general

import "net/url"

// GetUrlHost 解析 URL ，获取其主机名
//
// 参数：
//   - rawURL: 待解析的 URL
//
// 返回：
//   - URL 中的主机名
//   - 错误信息
func GetUrlHost(rawURL string) (string, error) {
	// 解析URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	return parsedURL.Host, nil

}
