package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Post发送json数据
// 返回响应内容、http响应码、错误信息
func PostJson(url string, data interface{}) (string, string, error) {
	d, err := json.Marshal(data)
	if err != nil {
		logs.Error(fmt.Sprintf("[PostJson] json数据转字符串失败 %s", err))
		return "", "", err
	}
	return PostJsonStr(url, string(d))
}

//Post发送json数据。返回响应内容、http响应码、错误信息
func PostJsonTimeout(url string, data interface{}, timeout time.Duration) (string, string, error) {
	d, err := json.Marshal(data)
	if err != nil {
		logs.Error(fmt.Sprintf("[PostJsonTimeout] json数据转字符串失败 %s", err))
		return "", "", err
	}
	return PostJsonStrTimeout(url, string(d), timeout)
}

// Post发送json字符串
// 返回响应内容、http响应码、错误信息
func PostJsonStr(url string, data string) (string, string, error) {
	createCertHttpResult, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		logs.Error(fmt.Sprintf("[PostJsonStr] 创建http请求失败 url地址%s %s", url, err))
		return "", "", err
	}
	createCertHttpResult.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(createCertHttpResult)
	if err != nil {
		logs.Error(fmt.Sprintf("[PostJsonStr] 获取请求%s响应失败 %s", url, err))
		return "", "", err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	return string(respBody), resp.Status, nil
}

// Post发送json字符串
// 返回响应内容、http响应码、错误信息
func PostJsonStrTimeout(url string, data string, timeout time.Duration) (string, string, error) {
	createCertHttpResult, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		logs.Error(fmt.Sprintf("[PostJsonStrTimeout] 创建http请求失败 url地址%s %s", url, err))
		return "", "", err
	}
	createCertHttpResult.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	client.Timeout = timeout
	resp, err := client.Do(createCertHttpResult)
	if err != nil {
		logs.Error(fmt.Sprintf("[PostJsonStrTimeout] 获取请求%s响应失败 %s", url, err))
		return "", "", err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	return string(respBody), resp.Status, nil
}

// 下载文件
func DownloadFileFromUrl(url string, output string) error {
	_, err := CreateFileIfNotExists(output)
	if err != nil {
		return err
	}
	out, err := os.OpenFile(output, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	defer out.Close()
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	pix, err := ioutil.ReadAll(resp.Body)
	_, err = io.Copy(out, bytes.NewReader(pix))
	return err
}
