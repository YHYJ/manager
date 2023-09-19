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
	"errors"
	"net/http"
	"time"
)

type Data struct {
	Name string `json:"name"`
}

func getLatestVersion(url string) (error, string) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err, ""
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err, ""
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status), ""
	}

	var datas []Data
	err = json.NewDecoder(resp.Body).Decode(&datas)
	if err != nil {
		return err, ""
	}

	latestVersion := datas[0].Name

	return nil, latestVersion
}
