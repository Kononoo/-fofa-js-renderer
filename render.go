package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"log"
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
	var body string
	var screen []byte

	// Listen for network events to capture response headers and status code
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventResponseReceived); ok {
			log.Printf("我的请求url: %s，response.url: %s", url, ev.Response.URL)
			if ev.Type == network.ResourceTypeDocument {
				result.StatusCode = int(ev.Response.Status)
				result.Header = http.Header{}
				for k, v := range ev.Response.Headers {
					result.Header.Set(k, fmt.Sprintf("%v", v))
				}
			}
		}
	})

	actions := []chromedp.Action{
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &body),
	}

	// 进行页面截屏
	if screenshot {
		actions = append(actions, chromedp.FullScreenshot(&screen, 90))
	}

	err := chromedp.Run(ctx, actions...)
	if err != nil {
		return result, fmt.Errorf("failed to render page: %w", err)
	}

	// 获取渲染时间
	endTime := time.Now()
	renderTime := endTime.Sub(start).Milliseconds()

	// 封装结果
	result.URL = url
	result.Body = body
	result.RenderTime = renderTime
	result.ScreenShot = screen

	return result, nil

}
