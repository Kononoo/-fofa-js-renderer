package main

import (
	"context"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

type ScanResult struct {
	URL          string      `json:"url"`
	StatusCode   int         `json:"status_code"`   // 首页状态码
	Body         string      `json:"body"`          // 渲染后的html->body
	Header       http.Header `json:"header"`        // 首页Header
	RenderTime   int64       `json:"render_time"`   // 渲染时间，单位毫秒
	ScreenShot   []byte      `json:"screenshot"`    // 截屏
	ErrorMessage string      `json:"error_message"` // 渲染失败、错误信息
}

func Scan(url string, screenshot bool) (ScanResult, error) {
	start := time.Now() // 记录渲染时间
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var result ScanResult // 渲染结果
	result.URL = url

	var body string
	var statusCode int
	var headers http.Header
	var screenshotBuf []byte

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &body),
		chromedp.Evaluate(`(() => {
			return fetch(window.location.href, {method: 'GET'}).then(response => {
				const headers = {};
				for (let [key, value] of response.headers.entries()) {
					headers[key] = value;
				}
				return { headers: headers, statusCode: response.status };
			});
		})()`, &headers),
	}

	// 进行页面截屏
	if screenshot {
		tasks = append(tasks, chromedp.FullScreenshot(&screenshotBuf, 90))
		result.ScreenShot = screenshotBuf
	}
	if err := chromedp.Run(ctx, tasks...); err != nil {
		result.ErrorMessage = err.Error()
		return result, err
	}

	// 封装结果
	result.Body = body
	result.Header = headers
	result.StatusCode = statusCode
	result.RenderTime = time.Since(start).Milliseconds()

	return result, nil

}
