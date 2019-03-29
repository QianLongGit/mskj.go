package vo

// 文件信息结构体
type FileInfo struct {
	// 文件名称
	Name string `json:"name"`
	// 文件绝对路径
	Abs string `json:"abs"`
	// 文件大小
	Size int64 `json:"size"`
	// 文件访问时间
	AccessTimestamp int64 `json:"access_time"`
	// 文件修改时间
	ModifyTimestamp int64 `json:"modify_time"`
	// 文件创建时间
	CreateTimestamp int64 `json:"create_time"`
}

