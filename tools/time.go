package tools

import (
	"github.com/astaxie/beego/logs"
	"time"
)

var golangBirthTime = "2006-01-02 15:04:05"
var golangBirthDate = "2006-01-02"

// 将当前时间转为yyyy-MM-dd hh:mm:ss格式
func GetNowTimeFormat() string {
	return GetNowTimeFormatByFormat(golangBirthTime)
}

// 将当前时间转为yyyy-MM-dd格式
func GetNowDateFormat() string {
	return GetNowTimeFormatByFormat(golangBirthDate)
}

// 指定时间戳, 获取其对应日期00:00:00点时间戳
func GetTimestampStartUnix(timestamp int64) int64 {
	tm := Unix2Time(timestamp)
	tmStr := GetTimeFormatByFormat(tm, golangBirthDate)
	tmStart := TimeParseByFormat(tmStr, golangBirthDate)
	return tmStart.Unix()
}

// 指定时间戳, 获取其对应日期23:59:59点时间戳
func GetTimestampEndUnix(timestamp int64) int64 {
	tm := Unix2Time(timestamp)
	tmStr := GetTimeFormatByFormat(tm, golangBirthDate)
	tmEnd := TimeParseByFormat(tmStr+" 23:59:59", golangBirthTime)
	return tmEnd.Unix()
}

// 时间戳转时间对象
func Unix2Time(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// 时间戳转time字符串
func Unix2TimeFormat(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(golangBirthTime)
}

// Now, 获取自定义format字符串
func GetNowTimeFormatByFormat(format string) string {
	return GetTimeFormatByFormat(time.Now(), format)
}

// 指定时间, 获取自定义format字符串
func GetTimeFormatByFormat(time time.Time, format string) string {
	return time.Format(format)
}

// 指定时间, 获取yyyy-MM-dd hh:mm:ss格式
func GetTimeFormat(time time.Time) string {
	return time.Format(golangBirthTime)
}

// 指定时间对象, 获取yyyy-MM-dd格式
func GetDateFormat(time time.Time) string {
	return time.Format(golangBirthDate)
}

// 指定时间字符串/字符串格式, 转换为时间对象
func TimeParseByFormat(timeStr string, format string) *time.Time {
	t, err := time.ParseInLocation(format, timeStr, time.Local)
	if err != nil {
		logs.Error(err)
		return nil
	}
	return &t
}
