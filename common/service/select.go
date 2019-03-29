package service

import (
	"fmt"
	"gitee.com/piupuer/go/tools"
	"gitee.com/piupuer/go/vo"
	"github.com/astaxie/beego/orm"
	"reflect"
	"strings"
)

// 数据库查询服务工具类
// 抽取了数据库操作中经常使用分页查询的一些方法

// 支持的filter列表
// 在struct中加入tag filter, 默认值exact
// 如果同时包含tag related, 将执行关联查询
// 其中related=true时, 按驼峰转换 如 MachineName string `filter:"contains" related:"true"`  =>  qs.Filter("machine__name__contains",value)
// 其中related=自定义值时, 按自定义转换 如 VasAppGpuId string `filter:"contains" related:"VasApp__GpuId"`  =>  qs.Filter("vasapp__gpuid__contains",value)
// 可结合官方文档学习：https://beego.me/docs/mvc/model/query.md#operators
var filterList = [...]string{
	// 适用于 LIKE '%...%'
	"contains",
	// 适用于 =
	"exact",
	// 适用于 >
	"gt",
	// 适用于 >=
	"gte",
	// 适用于 <
	"lt",
	// 适用于 <=
	"lte",
	// 适用于 IN (...)
	"in",
	// 适用于 LIKE '...%'
	"startswith",
	// 适用于 LIKE '%...'
	"endswith",
	// 适用于 IS NULL / IS NOT NULL
	"isnull",
}

// 查询器
type query struct {
	orm.QuerySeter
}

//============================================================
//=======================公共方法==============================
//============================================================
// 获取查询器
// pagination vo.Pagination	pagination(分页)
// params 为可变参数, 依次为
// 参数位置	类型				含义												必要性	示例
// [0] 		interface{}		entity(需要查询的实体)								是
// [1] 		interface{}		filter/cond										否		IdGt: 1(转为数据库语句WHERE id > 1)
//			注意filter和cond只能选择一个, 另一个将无法生效
// 			如果为filter表示使用Filter过滤条件, 支持的filter列表在该文件最上方
//			如果为cond表示使用Condition自定义条件
// [2] 		interface{}		order(排序条件)									否		Id: "asc"(转为数据库语句ORDER BY id asc)
func Query(pagination *vo.Pagination, params ...interface{}) orm.QuerySeter {
	// 创建查询
	q := query{}
	// 读数据使用从库
	o := LockOrmer(true)
	defer UnlockOrmer(true)
	// 获取可变参数的个数
	l := len(params)
	if l == 0 {
		panic("缺少entity实体参数, 无法继续查询")
	}
	// 1. 查询实体
	q.QuerySeter = o.QueryTable(params[0])
	if l >= 2 && params[1] != nil {
		// 2. 过滤条件
		// 类型断言: 如果类型为orm.Condition, 则视为cond, 否则为filter
		cond, err := params[1].(*orm.Condition)
		if err {
			// 添加condition
			q.QuerySeter = q.SetCond(cond)
		} else {
			// 添加过滤器
			q = q.filter(params[1])
		}
	}
	if l >= 4 && params[3] != nil {

	}
	// 3. 查询数据总条数
	total, err := q.Count()
	if err != nil {
		total = 0
	}
	if l >= 3 && params[2] != nil {
		// 4. 添加排序
		q = q.order(params[2])
	}
	// 5. 添加分页
	q = q.page(*pagination)
	// 赋值给pagination中total
	pagination.Total = total
	// 最后, 返回查询器
	return q.QuerySeter
}

//============================================================
//======================公共方法结束============================
//============================================================

//============================================================
//=======================私有方法==============================
//============================================================
// 添加过滤器
func (q query) filter(filters interface{}) query {
	filterType := reflect.TypeOf(filters)
	filterValue := reflect.ValueOf(filters)
	// 不允许为空值
	// 遍历所有属性
	for i := 0; i < filterType.NumField(); i++ {
		// 字段
		field := filterType.Field(i)
		// 获取filter标签
		filterTag := strings.ToLower(strings.TrimSpace(filterType.Field(i).Tag.Get("filter")))
		// 默认使用等于
		if filterTag == "" {
			filterTag = "exact"
		}
		// 获取related标签
		relatedTag := strings.TrimSpace(filterType.Field(i).Tag.Get("related"))
		// 属性key, 默认按驼峰处理
		key := field.Name
		if relatedTag != "" {
			if relatedTag == "true" {
				// 默认按驼峰转换, 转为双下划线, 符合beego要求
				key = tools.CamelStr2FormatStr(field.Name, "__")
			} else {
				// 设置为和tag一致的字符
				key = strings.ToLower(relatedTag)
			}
		}
		// 不允许为空值
		if isBlank(filterValue.Field(i)) {
			continue
		}
		// 属性值
		value := filterValue.Field(i).Interface()
		// 从支持的过滤列表中遍历
		for _, v := range filterList {
			// 过滤tag
			if filterTag == v {
				// 查询字段
				if fmt.Sprint(value) == "TRue" {
					//若字段值为TRue，则转换为bool值
					q.QuerySeter = q.Filter(fmt.Sprintf("%s__%s", key, filterTag), true)
				} else if fmt.Sprint(value) == "FAlse" {
					//若字段值为FAlse，则转换为bool值
					q.QuerySeter = q.Filter(fmt.Sprintf("%s__%s", key, filterTag), false)
				} else {
					q.QuerySeter = q.Filter(fmt.Sprintf("%s__%s", key, filterTag), value)
				}
				break
			}
		}
	}
	return q
}

// 添加排序
func (q query) order(orders interface{}) query {
	orderType := reflect.TypeOf(orders)
	orderValue := reflect.ValueOf(orders)
	// 不允许为空值
	if !isBlank(orderValue) {
		var arr []string
		// 遍历所有属性
		for i := 0; i < orderType.NumField(); i++ {
			// 字段
			field := orderType.Field(i)
			// 属性key
			key := strings.ToLower(field.Name)
			// 属性值
			value := orderValue.Field(i).Interface()
			if value == "asc" {
				// 顺序排列
				arr = append(arr, key)
			} else if value == "desc" {
				// 降序排列
				arr = append(arr, fmt.Sprintf("-%s", key))
			}
		}
		// 设置排序
		q.QuerySeter = q.OrderBy(arr...)
	}
	return q
}

// 获取带分页的查询器
// 由于go不支持泛型，因此不能直接获取结果集
func (q query) page(pagination vo.Pagination) query {
	// 加入分页限制
	q.QuerySeter = q.Offset((pagination.Start - 1) * pagination.Limit).Limit(pagination.Limit)
	return q
}

// 判断value是否空值
func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		// 查询时bool的值是有针对性的, 因此这里返回false
		return false
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}

//============================================================
//======================私有方法结束============================
//============================================================
