package tools

import (
	"regexp"
	"strings"
)

// 字符串工具

// 驼峰转下划线, XxYy to xx_yy , XxYY to xx_yy
func CamelStr2UnderlineStr(s string) string {
	return CamelStr2FormatStr(s, "_")
}

// 下划线转驼峰, xx_yy to XxYy
func UnderlineStr2CamelStr(s string) string {
	return FormatStr2CamelStr(s, "_")
}

// 驼峰转固定格式字符串, XxYy to xx_yy , XxYY to xx_yy
func CamelStr2FormatStr(s string, format string) string {
	// 只匹配大写字母
	pat := "[A-Z]"
	regexp.Match(pat, []byte(s))
	re, _ := regexp.Compile(pat)
	// 将匹配到的部分,将字符转为小写, 在前面添加指定格式字符串
	str := re.ReplaceAllStringFunc(s, func(b string) string {
		return format + strings.ToLower(b)
	})
	// 去掉首个format
	return strings.TrimLeft(str, format)
}

// 固定格式转驼峰字符串, xx{format}yy to XxYy
func FormatStr2CamelStr(s string, format string) string {
	// 匹配format+单个小写字母
	re := regexp.MustCompile(format + "[a-z]")
	// 将匹配到的部分,将字符转为小写, 在前面添加指定格式字符串
	str := re.ReplaceAllStringFunc(s, func(b string) string {
		// 去掉format, 转为大写
		return strings.ToUpper(strings.TrimLeft(b, format))
	})
	return str
}

// 去除多余的文件分隔符只保留一个
func RemoveFileSeparator(s string) string {
	n := strings.Replace(s, "//", "/", -1)
	// 去除前后没有变化, 则直接返回, 否则递归
	if n == s {
		return s
	}
	return RemoveFileSeparator(n)
}
