package tools

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

// 类型转换

// struct转为包含相同字段的struct结构体, 中间通过json桥接
func Struct2StructByJson(m interface{}, s interface{}) error {
	machineJson, err := json.Marshal(m)
	if err != nil {
		logs.Error(fmt.Sprintf("map转换json失败 异常 %s", err))
		return err
	}

	err = json.Unmarshal(machineJson, s)
	if err != nil {
		logs.Error(fmt.Sprintf("json转换map失败 异常 %s", err))
		return err
	}
	return nil
}
