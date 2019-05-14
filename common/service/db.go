package service

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"mskj.go/common"
	"sync"
	"sync/atomic"
)

// sqlite不支持多个事务同时执行, 这里对db操作重构

// 数据库锁
var dbLock = new(sync.Mutex)

// 连接请求数
var dbCount int64
var dbReadCount int64
var dbWriteCount int64

// 获取Orm锁
func LockOrmer(read bool) orm.Ormer {
	atomic.AddInt64(&dbCount, 1)
	if read {
		atomic.AddInt64(&dbReadCount, 1)
	} else {
		atomic.AddInt64(&dbWriteCount, 1)
	}
	dbLock.Lock()
	if read {
		return common.GetSlaveNewOrm()
	}
	return common.GetNewOrm()
}

// 释放Orm锁
func UnlockOrmer(read bool) {
	dbLock.Unlock()
	atomic.AddInt64(&dbCount, -1)
	if read {
		atomic.AddInt64(&dbReadCount, -1)
	} else {
		atomic.AddInt64(&dbWriteCount, -1)
	}
	if atomic.LoadInt64(&dbCount) < 0 {
		logs.Warn("Orm锁连接请求数小于0, 请检查锁释放是否有遗漏")
	}
}

// 获取连接数
func GetDbCount() int64 {
	return atomic.LoadInt64(&dbCount)
}

func GetDbReadCount() int64 {
	return atomic.LoadInt64(&dbReadCount)
}

func GetDbWriteCount() int64 {
	return atomic.LoadInt64(&dbWriteCount)
}
