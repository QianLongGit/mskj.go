package vo

// 分页包装对象
type Pagination struct {
	Start int   `json:"start" form:"pagination[start]"`
	Limit int   `json:"limit" form:"pagination[limit]"`
	Total int64 `json:"total" form:"pagination[total]"`
}
