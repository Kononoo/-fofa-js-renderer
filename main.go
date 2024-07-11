package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	render(urls)
}

// 页面渲染、保存结果
func render(urls []string) {
	for _, url := range urls {
		result, err := Scan(url, screenshot)
		if err != nil {
			log.Printf("Error scanning %s: %v\n", url, err)
			continue
		}

		resultJson, err := json.Marshal(result)
		if err != nil {
			log.Printf("Error marshalling json: %s %v\n", url, err)
			continue
		}

		resultFile := fmt.Sprintf("%s.json", strings.ReplaceAll(url, "https://", ""))
		if err = ioutil.WriteFile(resultFile, resultJson, 0644); err != nil {
			log.Printf("Error writing result to file for %s: %v\n", resultFile, err)
			continue
		}

		if screenshot {
			screenshotFile := fmt.Sprintf("%s.png", strings.ReplaceAll(url, "https://", ""))
			if err := ioutil.WriteFile(screenshotFile, resultJson, 0644); err != nil {
				log.Printf("Error writing result to file for %s: %v\n", screenshotFile, err)
				continue
			}

		}

	}
}
