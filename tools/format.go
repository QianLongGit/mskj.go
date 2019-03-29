package tools

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// 根据日期格式化字符串, 获取文件大小(B)
func GetFileSizeByFormat(format string) int64 {
	// 格式化
	arr := []string{
		"M",
		"G",
		"T",
		"P",
		"E",
		"Z",
	}
	var num int64 = 0
	i := 0
	format = strings.ToUpper(format)
	for i = 0; i < len(arr); i++ {
		if strings.HasSuffix(format, arr[i]) {
			format = strings.Replace(format, arr[i], "", -1)
			num, _ = strconv.ParseInt(format, 10, 64)
			break
		}
	}

	// 不同级别对应不同字符串
	switch i {
	case 0:
		return num * 1024 * 1024
	case 1:
		return num * 1024 * 1024 * 1024
	case 2:
		return num * 1024 * 1024 * 1024 * 1024
	case 3:
		return num * 1024 * 1024 * 1024 * 1024
	case 4:
		return num * 1024 * 1024 * 1024 * 1024 * 1024
	case 5:
		return num * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return num
	}
}

// 根据日期格式化字符串, 获取时间戳
func GetSecondByFormat(format string) int64 {
	// 格式化
	arr := []string{
		// 秒
		"s",
		// 分钟
		"m",
		// 小时
		"h",
		// 天
		"d",
		// 月
		"M",
		// 年
		"y",
	}
	var num int64 = 0
	i := 0
	format = strings.ToLower(format)
	for i = 0; i < len(arr); i++ {
		if strings.HasSuffix(format, arr[i]) {
			format = strings.Replace(format, arr[i], "", -1)
			num, _ = strconv.ParseInt(format, 10, 64)
			break
		}
	}

	// 不同级别对应不同字符串
	switch i {
	case 0:
		return num * 1000
	case 1:
		return num * 60000
	case 2:
		return num * 3600000
	case 3:
		return num * 86400000
	case 4:
		return num * 2592000000
	case 5:
		return num * 31104000000
	default:
		return num
	}
}

// 获取文件格式化后缀
func FileSizeSuffix(level int) string {
	// 不同级别对应不同字符串
	switch level {
	case 1:
		return "B"
	case 2:
		return "K"
	case 3:
		return "M"
	case 4:
		return "G"
	case 5:
		return "T"
	case 6:
		return "P"
	case 7:
		return "E"
	case 8:
		return "Z"
	case 9:
		return "Y"
	default:
		return "B"
	}
}

// 自动获取临界线, 该除以多少个1024
func getCritical(size int64) (int, int64) {
	var level = 1
	num := float64(size)
	// 临界线
	var critical int64 = 1
	for {
		if num < 1024 || level >= 9 {
			break
		}
		num = num / 1024
		critical *= 1024
		level ++
	}
	return level, critical
}

// 文件大小格式化 B/KB/MB/GB/TB/PB/EB/ZB/YB
func FileSizeFormat(size int64) string {
	_, format := FileSizeFormatByLevel(size, 0)
	return format
}

// 文件大小格式化, 指定等级 B/KB/MB/GB/TB/PB/EB/ZB/YB
// 返回值1 不带B/KB末尾字符 返回值2 带B/KB末尾字符
func FileSizeFormatByLevel(size int64, level int) (string, string) {
	if level > 0 {
		// 指定级别, 手动获取临界值
		critical := math.Pow(1024, float64(level-1))
		return fmt.Sprintf("%.2f", float64(size)/critical), fmt.Sprintf("%.3f%s", float64(size)/critical, FileSizeSuffix(level))
	} else {
		// 不指定级别, 自动获取临界值
		level, critical := getCritical(size)
		return fmt.Sprintf("%.2f", float64(size)/float64(critical)), fmt.Sprintf("%.3f%s", float64(size)/float64(critical), FileSizeSuffix(level))
	}
}

// 文件大小格式化为同一个级别, 以第一个为基准
func FileSizeSameLevelFormat(sizes ...int64) []string {
	// 获取第一个值的级别
	level, _ := getCritical(sizes[0])
	suffix := FileSizeSuffix(level)
	var res []string
	for i := 0; i < len(sizes); i++ {
		size, _ := FileSizeFormatByLevel(sizes[i], level)
		res = append(res, size)
	}
	// 将level和后缀加入到结果集
	res = append(res, strconv.Itoa(level))
	res = append(res, suffix)
	return res
}

// 时间戳格式化 秒/分钟/小时/天
func TimestampFormat(timestamp int64) string {
	if timestamp < 60000 {
		return fmt.Sprintf("%.2f秒", float64(timestamp)/1000)
	}
	if timestamp < 3600000 {
		return fmt.Sprintf("%.2f分", float64(timestamp)/60000)
	}
	if timestamp < 86400000 {
		return fmt.Sprintf("%.2f时", float64(timestamp)/3600000)
	}
	if timestamp < 2592000000 {
		return fmt.Sprintf("%.2f天", float64(timestamp)/86400000)
	}
	return fmt.Sprintf("%.2f月", float64(timestamp)/2592000000)
}

// 百分比格式化
func PercentageFormat(dividend int64, divisor int64) string {
	if dividend == 0 || divisor == 0 {
		return "0.000%"
	}
	// 输出%本身
	return fmt.Sprintf("%.3f%%", float64(dividend)/float64(divisor)*100)
}
