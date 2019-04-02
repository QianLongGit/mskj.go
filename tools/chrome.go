package tools

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"log"
	"time"
)

// 爬虫工具, 主要使用chromedp/chromedp, 后台调用chrome浏览器

// 借助chromeless访问指定页面(适合网页中存在异步数据的场景)
func VisitUrlWithChromeLess(url string, ele string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cdp, err := newHeadless(ctx)
	html := ""
	if err != nil {
		return html, errors.New(fmt.Sprintf("创建Chrome Headless实例失败 异常 %s", err.Error()))
	}
	if ele == "" {
		ele = "body"
	}
	// cdp是chromedp实例
	// ctx是创建cdp时使用的context.Context
	err = cdp.Run(ctx, chromedp.Tasks{
		// 访问指定URL
		chromedp.Navigate(url),
		// 等待直到页面元素加载完毕
		chromedp.WaitVisible(ele, chromedp.ByQuery),
		// 获取HTML
		chromedp.OuterHTML(ele, &html, chromedp.ByQuery),
	})
	if err != nil {
		return html, errors.New(fmt.Sprintf("页面渲染失败 异常 %s", err.Error()))
	}
	return html, nil
}

// NewHeadless 创建headless chrome实例
// chromedp内部有自己的超时设置，你也可以通过ctx来设置更短的超时
func newHeadless(ctx context.Context) (*chromedp.CDP, error) {
	// runner.Flag设置启动headless chrome时的命令行参数
	// runner.URL设置启动时打开的URL
	// Windows用户需要设置runner.Flag("disable-gpu", true)，具体信息参见文档的FAQ
	run, err := runner.New(runner.Flag("headless", true))

	if err != nil {
		return nil, err
	}

	// run.Start启动实例
	err = run.Start(ctx)
	if err != nil {
		return nil, err
	}

	// 默认情况chromedp会输出大量log，因为是示例所以选择屏蔽，dropChromeLogs为自定义函数，形式为func(string, ...interface{}){}
	// 使用runner初始化chromedp实例
	// 实例在使用完毕后需要调用c.Shutdown()来释放资源
	c, err := chromedp.New(ctx, chromedp.WithRunner(run), chromedp.WithErrorf(log.Printf))
	if err != nil {
		return nil, err
	}

	return c, nil
}
