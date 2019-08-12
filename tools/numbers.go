package tools

import (
	"github.com/astaxie/beego/logs"
	"strconv"
)

//字符串转int64类型，转换错误将会输出日志并返回0
func StringToInt64(number string) int64 {
	if number == "" {
		return 0
	}
	i, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		logs.Error(err)
		return 0
	}
	return i
}

//字符串转int64类型，转换错误将会输出日志并返回0
func StringToInt(number string) int {
	if number == "" {
		return 0
	}
	i, err := strconv.Atoi(number)
	if err != nil {
		logs.Error(err)
		return 0
	}
	return i
}

//字符串转float64类型，转换错误将会输出日志并返回0
func StringToFloat64(number string) float64 {
	if number == "" {
		return 0
	}
	i, err := strconv.ParseFloat(number, 64)
	if err != nil {
		logs.Error(err)
		return 0
	}
	return i
}
