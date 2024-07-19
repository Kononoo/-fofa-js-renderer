package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	url        string
	file       string
	query      string
	screenshot bool
)

// 命令行工具-启动...
func main() {
	flag.StringVar(&url, "u", "", "The URL to scan")
	flag.StringVar(&file, "f", "", "The file containing urls to scan")
	flag.StringVar(&query, "q", "", "The FOFA query to search")
	flag.BoolVar(&screenshot, "s", false, "Whether to take screenshot	")
	flag.Parse()

	if url == "" && file == "" && query == "" {
		flag.Usage()
		log.Fatal("You must provide either a URL or a file containing URLs")
	}

	// 所有要渲染的网站url
	var urls []string
	if url != "" {
		urls = []string{url}
	}

	// 读取文件中的url
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal("文件打开失败", err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	// 获取配置的apiKey
	config, err := loadConfig()
	log.Println("成功加载配置文件...")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 读取查询参数，进行FOFA-API调用查询url
	if query != "" {
		fofaUrls, err := getFofaResults(config.ApiKey, query, config.MaxUrls)
		if err != nil {
			log.Fatalf("Error fetching FOFA results: %v", err)
		}
		urls = append(urls, fofaUrls...)
	}

	// 创建result目录
	if err := os.MkdirAll("result", os.ModePerm); err != nil {
		log.Fatalf("Error creating result directory: %v", err)
	}

	renderFunc(urls)
}

// 页面渲染、保存结果
func renderFunc(urls []string) {
	var wg sync.WaitGroup
	concurrency := 50                           // 设置并发数量
	sem := make(chan struct{}, concurrency)     // 控制并发数的信号量
	results := make(chan ScanResult, len(urls)) // 存储渲染结果通道

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			sem <- struct{}{} // 获取一个信号量
			result, err := Scan(url, screenshot)
			if err != nil {
				result.ErrorMessage = err.Error()
			}
			results <- result
			<-sem // 释放一个信号量
		}(url)
	}

	wg.Wait()

	// 关闭结果通道
	go func() {
		close(results)
	}()

	// 保存渲染结果
	for result := range results {
		if result.URL != "" {
			log.Printf("Skipping empty result for URL: %s", result.URL)
		} else {
			saveResult(result)
		}
	}
}

// 保存结果到文件
func saveResult(result ScanResult) {
	resultJson, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error marshalling json: %s %v\n", result.URL, err)
		return
	}

	urlFileSafe := createSafeFileName(result.URL)

	// 保存渲染结果JSON文件
	resultFile := filepath.Join("result", fmt.Sprintf("%s.json", urlFileSafe))
	if err := ioutil.WriteFile(resultFile, resultJson, 0644); err != nil {
		log.Printf("Error writing result to file for %s: %v", result.URL, err)
		return
	}
	log.Printf("成功保存文件->JSON：%s", resultFile)

	// 如果需要截图则保存截图文件
	if screenshot {
		screenshotFile := filepath.Join("result", fmt.Sprintf("%s.png", urlFileSafe))
		if err := ioutil.WriteFile(screenshotFile, result.ScreenShot, 0644); err != nil {
			log.Printf("Error writing screenshot to file for %s: %v", result.URL, err)
		} else {
			log.Printf("成功保存文件->图片：%s", screenshotFile)
		}
	}
}

// 生成安全的文件名
func createSafeFileName(url string) string {
	urlFileSafe := strings.ReplaceAll(url, "https://", "")
	urlFileSafe = strings.ReplaceAll(urlFileSafe, "http://", "")
	urlFileSafe = strings.ReplaceAll(urlFileSafe, "/", "_")
	urlFileSafe = strings.ReplaceAll(urlFileSafe, "?", "_")
	urlFileSafe = strings.ReplaceAll(urlFileSafe, "&", "_")
	urlFileSafe = strings.ReplaceAll(urlFileSafe, "=", "_")
	urlFileSafe = strings.ReplaceAll(urlFileSafe, ":", "_")
	if len(urlFileSafe) > 100 {
		urlFileSafe = urlFileSafe[:100] // 限制文件名长度，避免文件名过长
	}
	return urlFileSafe
}
