package tools

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// 随机数字生成器: 指定数字位数、生成个数
func GenRandomNums(d int, l int) ([]int, error) {
	// 剩余配额
	remaining := int(math.Pow10(d) - math.Pow10(d-1))

	if d == 1 {
		remaining = 10
	}
	if l == -1 {
		l = remaining
	} else if l < -1 {
		return nil, fmt.Errorf("数据个数不对")
	}

	if remaining < l {
		return nil, fmt.Errorf("剩余配额不足")
	}

	return GenRandomNumsExceptNums(d, l, []int{})
}

// 随机数字生成器: 指定数字位数、生成个数、排除指定数组
func GenRandomNumsExceptNums(d int, l int, e []int) ([]int, error) {
	all, err := GenNumsExceptNums(d, l, e)
	// 结果集
	var res []int
	if err != nil {
		return res, err
	}
	for {
		rand.Seed(time.Now().Unix())
		// 随机获取索引
		index := rand.Intn(len(all))
		// 保存到结果集
		res = append(res, all[index])
		// 重置基准数组
		all = append(all[:index], all[index+1:]...)
		if len(res) == l {
			break
		}
	}
	return res, nil
}

// 数字生成器: 指定数字位数、生成个数
func GenNums(d int, l int) ([]int, error) {
	// 剩余配额
	remaining := int(math.Pow10(d) - math.Pow10(d-1))

	if d == 1 {
		remaining = 10
	}
	if l == -1 {
		l = remaining
	} else if l < -1 {
		return nil, fmt.Errorf("数据个数不对")
	}

	if remaining < l {
		return nil, fmt.Errorf("剩余配额不足")
	}

	return GenNumsExceptNums(d, l, []int{})
}

// 数字生成器: 指定数字位数、生成个数、排除指定数组
func GenNumsExceptNums(d int, l int, e []int) ([]int, error) {
	end := int(math.Pow10(d))
	start := int(math.Pow10(d - 1))
	if d == 1 {
		start = 0
	}
	// 剩余配额
	remaining := end - start
	if l == -1 {
		l = remaining - len(e)
	} else if l < -1 {
		return nil, fmt.Errorf("数据个数不对")
	}

	if remaining-len(e) < l {
		return nil, fmt.Errorf("剩余配额不足")
	}

	// 获取基准数组
	var all []int
	for i := start; i < end; i++ {
		if IsExistIntArray(e, i) == -1 {
			all = append(all, i)
		}
	}
	return all, nil
}
