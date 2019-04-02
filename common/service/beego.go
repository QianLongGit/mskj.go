package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
)

// 读取实体对象所有的关联实体
func LoadAllRelated(model interface{}) {
	// 读数据使用从库
	o := LockOrmer(true)
	defer UnlockOrmer(true)
	pType := reflect.TypeOf(model).Elem()
	for k := 0; k < pType.NumField(); k++ {
		ormTag := pType.Field(k).Tag.Get("orm")
		if strings.Contains(ormTag, "rel") || strings.Contains(ormTag, "reverse") {
			_, err := o.LoadRelated(model, pType.Field(k).Name)
			if err != nil && err != orm.ErrNoRows {
				logs.Error("Load Related Error", err)
			}
		}
	}
}
