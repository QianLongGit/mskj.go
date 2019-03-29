package tools

// 整数类型求绝对值
// 由于math.abs不支持整型, 如果转为float再转为int64耗时较多
func AbsInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
