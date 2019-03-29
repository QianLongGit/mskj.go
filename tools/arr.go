package tools

// 是否在数组中
func IsExistStringArray(arr []string, value string) int {
	for index, v := range arr {
		if v == value {
			return index
		}
	}
	return -1
}

func IsExistInt64Array(array []int64, value int64) int {
	for index, v := range array {
		if v == value {
			return index
		}
	}
	return -1
}

func IsExistIntArray(array []int, value int) int {
	for index, v := range array {
		if v == value {
			return index
		}
	}
	return -1
}
