package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FofaResponse struct {
	Results [][]string `json:"results"`
}

func getFofaResults(apiKey, query string) ([]string, error) {
	// 获取配置的apiKey
	config, _ := loadConfig()
	apiKey = config.ApiKey

	// 对请求query进行base64编码
	encodedQuery := base64.StdEncoding.EncodeToString([]byte(query))
	url := fmt.Sprintf("https://fofa.info/api/v1/search/all?key=%s&qbase64=%s&fields=host&size=3000", apiKey, encodedQuery)

	req, err := http.NewRequest("Get", url, nil)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get fofa results -> status code: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var fofaResponse FofaResponse
	if err = json.Unmarshal(body, &fofaResponse); err != nil {
		return nil, err
	}

	var urls []string
	for _, result := range fofaResponse.Results {
		urls = append(urls, result[0])
	}

	return urls, nil
}
