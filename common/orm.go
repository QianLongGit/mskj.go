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
	err := o.Using("slave")
	if err != nil {
		o = orm.NewOrm()
	}
	return o
}
