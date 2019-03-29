package tools

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

// 获取单核CPU平均系统负载
func LoadAvgPerCpu() (float64, error) {
	avg, err := load.Avg()
	if err != nil {
		return 0, err
	}
	counts, err := cpu.Counts(false)
	if err != nil {
		return 0, err
	}
	return avg.Load1 / float64(counts), nil
}

// 判断系统是否超负荷运转(平均负载是否超出指定值)
func IsOverLoad(num float64) (float64, bool) {
	avg, err := LoadAvgPerCpu()
	if err == nil {
		if avg > num {
			return avg, true
		}
	}
	return avg, false
}
