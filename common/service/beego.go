package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
)

// 读取实体对象所有的关联实体
// 谨慎使用该方法，该方法会造成性能影响。推荐使用 LoadModelRelated 方法，指定具体的关联属性查询
func LoadAllRelated(model interface{}) {
	pType := reflect.TypeOf(model).Elem()
	var cols []string
	for k := 0; k < pType.NumField(); k++ {
		ormTag := pType.Field(k).Tag.Get("orm")
		if strings.Contains(ormTag, "rel") || strings.Contains(ormTag, "reverse") {
			col := pType.Field(k).Name
			cols = append(cols, col)
		}
	}
	// 获取cols关联实体
	LoadModelRelated(model, cols...)
}

// 获取指定列关联实体
func LoadModelRelated(model interface{}, cols ...string) {
	// 读数据使用从库
	o := LockOrmer(true)
	defer UnlockOrmer(true)
	for _, col := range cols {
		_, err := o.LoadRelated(model, col)
		if err != nil && err != orm.ErrNoRows {
			logs.Error(fmt.Sprintf("Load Related %s Error %s", col, err))
		}
	}
}
