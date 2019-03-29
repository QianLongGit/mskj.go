package tools

import (
	"fmt"
	"testing"
)

// 获取系统负载
func TestLoadAvgPerCpu(t *testing.T) {
	fmt.Println(LoadAvgPerCpu())
}
