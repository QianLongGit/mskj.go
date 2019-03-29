package common

import (
	"github.com/astaxie/beego/orm"
)

// beego的数据库操作对象
// 主库 write
func GetNewOrm() orm.Ormer {
	return orm.NewOrm()
}

// 从库 read
func GetSlaveNewOrm() orm.Ormer {
	o := orm.NewOrm()
	// 这里暂时和主库使用同一个数据库
	// o.Using("slave")
	return o
}
