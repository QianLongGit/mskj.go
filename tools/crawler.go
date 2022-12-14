package tools

import (
	"github.com/QianLongGit/mskj.go/s_const"
	"github.com/gocolly/colly"
)

// 爬虫实例

// 基于colly获取爬虫实例
func CollyCrawler(url string, response colly.ResponseCallback) error {
	// 创建爬虫实例
	c := colly.NewCollector()

	// 模拟浏览器行为
	c.UserAgent = s_const.USER_AGENT

	// 请求之前
	c.OnRequest(setHeaders)

	// 请求响应
	c.OnResponse(func(r *colly.Response) {
		response(r)
	})

	// 开始访问
	err := c.Visit(url)
	return err
}

// 基于colly获取爬虫实例
func CollyCrawlerWithHTML(url string, selector string, html colly.HTMLCallback) error {
	// 创建爬虫实例
	c := colly.NewCollector()

	// 模拟浏览器行为
	c.UserAgent = s_const.USER_AGENT

	// 请求之前
	c.OnRequest(setHeaders)

	// 请求响应
	c.OnHTML(selector, func(elem *colly.HTMLElement) {
		html(elem)
	})

	// 开始访问
	err := c.Visit(url)
	return err
}

// 设置请求头
func setHeaders(r *colly.Request) {
	r.Headers.Set("Referer", "https://www.baidu.com")
	r.Headers.Set("Accept-Encoding", "gzip, deflate")
	r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	r.Headers.Set("Connection", "keep-alive")
	r.Headers.Set("Accept", "*/*")
}
