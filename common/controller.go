package common

import (
	"encoding/json"
	"fmt"
	"github.com/QianLongGit/mskj.go/vo"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

// 针对beego的controller中没有的一些方法进行扩展

type Controller struct {
	beego.Controller
}

// 带消息的结果集
type RespMsg struct {
	// 状态码
	Code int `json:"code"`
	// 消息提示
	Msg string `json:"msg"`
}

// 带数据的结果集
type RespData struct {
	// 状态码
	Code int `json:"code"`
	// 结果集
	Data interface{} `json:"data"`
}

// 带分页数据的结果集
type RespPageData struct {
	RespData
	Pagination vo.Pagination `json:"pagination"`
}

// 指定状态码
func (controller *Controller) abort(code int, body string) {
	controller.CustomAbort(code, body)
}

// 自定义状态码, 同时回写错误信息
func (controller *Controller) msg(code int, msg string) {
	msgResp := RespMsg{
		Code: code,
		Msg:  msg,
	}
	// 转换JSON
	resp, _ := json.Marshal(msgResp)
	controller.abort(msgResp.Code, string(resp))
}

// 自定义状态码, 同时回写数据
func (controller *Controller) data(code int, data interface{}) {
	dataResp := RespData{
		Code: code,
		Data: data,
	}
	// 转换JSON
	resp, err := json.Marshal(dataResp)
	if err != nil {
		logs.Error(fmt.Errorf("转换json失败"))
		controller.Write500("转换json失败")
	}
	controller.abort(dataResp.Code, string(resp))
}

// 自定义状态码, 同时回写分页数据
func (controller *Controller) pageData(code int, data interface{}, pagination vo.Pagination) {
	pageDataResp := RespPageData{
		Pagination: pagination,
	}
	pageDataResp.Code = code
	pageDataResp.Data = data
	// 转换JSON
	resp, err := json.Marshal(pageDataResp)
	if err != nil {
		logs.Error(fmt.Errorf("转换json失败"))
		controller.Write500("转换json失败")
	}
	controller.abort(pageDataResp.Code, string(resp))
}

// 状态码200
func (controller *Controller) Write() {
	controller.msg(200, "请求成功")
}

// 状态码200, 同时回写数据
func (controller *Controller) WriteData(data interface{}) {
	controller.data(200, data)
}

// 状态码200, 同时回写分页数据
func (controller *Controller) WritePageData(data interface{}, pagination vo.Pagination) {
	controller.pageData(200, data, pagination)
}

// 状态码400, 同时回写错误信息
func (controller *Controller) Write400(msg string) {
	controller.msg(400, msg)
}

// 状态码500, 同时回写错误信息
func (controller *Controller) Write500(msg string) {
	controller.msg(500, msg)
}
