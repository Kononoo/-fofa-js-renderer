package main

import (
	"context"
	"fmt"
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

func Scan(url string, screenShot bool) (ScanResult, error) {
	start := time.Now() // 记录渲染时间
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var result ScanResult // 渲染结果
	result.URL = url

	var body string
	var headers http.Header
	var screenshotBuf []byte

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &body),
		chromedp.ActionFunc(func(ctx context.Context) error {
			res, exp, err := chromedp.Evaluate("document.readyState", nil).Do(ctx)
			if err != nil {
				return err
			}
			if res != "complete" {
				return fmt.Errorf("page not fully loaded：%s", exp)
			}
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			resp, err := chromedp.ExecAllocator(ctx).Target(ctx).Response()
			if err != nil {
				return err
			}
			headers = resp.Header
			result.StatusCode = resp.Status
			return nil
		}),
	}
	if screenshot {
		tasks = append(tasks, chromedp.FullScreenshot(&screenshotBuf, 90))
	}
	if err := chromedp.Run(ctx, tasks...); err != nil {
		result.ErrorMessage = err.Error()
		return result, err
	}

	result.Body = body
	result.Header = headers
	result.RenderTime = time.Since(start).Milliseconds()

	if screenshot {
		result.ScreenShot = screenshotBuf
	}

	return result, nil

}
