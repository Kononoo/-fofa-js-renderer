# fofa-js-renderer

## 项目简介
#### JS 渲染爬虫 —— Golang 实现 JS 渲染功能
目前的 FoFa Web 爬虫无法执行 JS 渲染，需要使用 Golang 实现 JS 渲染功能。
要求如下：
1. 渲染功能必须稳定，避免内存溢出（OOM），并能正确处理错误。
2. 确保对各种类型的 JS 页面具有良好的兼容性。
3. 支持使用本地 Chrome 和 Browserless 集群进行渲染。

## 使用标准
1. 使用 Golang 实现一个命令行工具，输入url或读取url列表文件进行渲染。保存每个url的结果为一个json文件；截屏保存为图片文件。
2. 必须能够准确识别并渲染 SPA 页面，确保数据完整获取。
3. 从fofa上取3000个目标，针对存活的网站能完成正确的渲染，确保正确性，稳定性。
4. 实现一个函数 Scan(url string, screenShot bool) (ScanResult, error)执行渲染动作。
```
type ScanResult struct {
    URL string
    StatusCode int // 首页状态码
    Body string // 渲染后的html body
    Header string // 首页的header 
    RenderTime int64 // 渲染时间，单位为毫秒
    ScreenShot []byte // 截屏
    RenderTime int64 // 渲染时间，单位为毫秒·
    ErrorMessage string // 如果渲染失败，记录错误信息
}
```

## 功能说明
- -u：单个 URL
> 描述：指定单个 URL 进行渲染。
示例：-u "https://example.com"
- -f：文件路径
> 描述：指定一个文件，文件中包含多个 URL，每行一个 URL。
示例：-f "urls.txt"
- -q：FOFA 查询语句
> 描述：指定一个 FOFA 查询语句，通过 FOFA API 获取 URL 列表进行渲染。
示例：-q "title=baidu.com"
- -s：是否截图 
> 描述：指定是否对页面进行截图，生成截图文件。
示例：-s


## 代码说明
- config.yml：配置文件，存放 FOFA API Key 和最大 URL 数量。
- main.go：主入口文件，处理命令行参数，调用渲染函数。
- render.go：包含 Scan 函数，执行 JS 渲染。
- config.go：包含加载配置文件的函数。
- result/：保存渲染结果的目录。



## 学习目标
- 了解 JS 渲染爬虫的概念和作用。
- 学习如何使用 Golang 实现简单的 JS 渲染爬虫工具。
- 理解 SPA 页面和传统页面的区别，以及其渲染方式。
- 探索不同类型的 Web 页面及其特征。

## 参考资料
- 编程语言文档和教程：学习 Golang 基础语法和编程概念。
- 网络爬虫相关技术：学习使用 Golang 实现网络爬虫的技术和方法。
- 开源工具和示例代码：搜索互联网上的开源 JS 渲染爬虫工具，学习其实现原理和代码结构，可以作为参考和学习的资源。
