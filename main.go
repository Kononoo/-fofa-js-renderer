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
)

var (
	url        string
	file       string
	screenshot bool
)

// 命令行工具-启动...
func main() {
	flag.StringVar(&url, "u", "", "The URL to scan")
	flag.StringVar(&file, "f", "", "The file containing urls to scan")
	flag.BoolVar(&screenshot, "s", false, "Whether to take screenshot	")
	flag.Parse()

	if url == "" && file == "" {
		flag.Usage()
		log.Fatal("You must provide either a URL or a file containing URLs")
	}

	// 所有要渲染的网站url
	var urls []string
	if url != "" {
		urls = []string{url}
	}
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

	// 创建result目录
	if err := os.MkdirAll("result", os.ModePerm); err != nil {
		log.Fatalf("Error creating result directory: %v", err)
	}

	renderFunc(urls)
}

// 页面渲染、保存结果
func renderFunc(urls []string) {
	for _, url := range urls {
		result, err := Scan(url, screenshot)
		if err != nil {
			log.Printf("Error scanning %s: %v\n", url, err)
			continue
		}

		resultJSON, err := json.Marshal(result)
		if err != nil {
			log.Printf("Error marshalling json: %s %v\n", url, err)
			continue
		}

		urlFileSafe := strings.ReplaceAll(url, "https://", "")
		urlFileSafe = strings.ReplaceAll(urlFileSafe, "http://", "")
		urlFileSafe = strings.ReplaceAll(urlFileSafe, "/", "_")

		resultFile := filepath.Join("result", fmt.Sprintf("%s.json", urlFileSafe))
		if err := ioutil.WriteFile(resultFile, resultJSON, 0644); err != nil {
			log.Printf("Error writing result to file for %s: %v", url, err)
			continue
		}

		if screenshot {
			screenshotFile := filepath.Join("result", fmt.Sprintf("%s.png", urlFileSafe))
			if err := ioutil.WriteFile(screenshotFile, result.ScreenShot, 0644); err != nil {
				log.Printf("Error writing screenshot to file for %s: %v", url, err)
				continue
			}
		}

	}
}
