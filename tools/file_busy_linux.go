package tools

import (
	"gitee.com/piupuer/go/s_const"
	"gitee.com/piupuer/go/vo"
	"os"
	"syscall"
	"time"
)

// 运行时: Linux

// 判断文件是否被占用
func IsFileBusyByInterval(filename string, busyInterval int64) (bool, *vo.FileInfo) {
	// 只能在1毫秒到1小时之间
	if busyInterval < 0 || busyInterval > 360000 {
		busyInterval = s_const.BUSY_INTERVAL
	}
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		// 无权限访问文件
		return true, nil
	}
	fi, err := f.Stat()
	if err != nil {
		// 无权限访问文件
		return true, nil
	}
	fileInfo := vo.FileInfo{}
	fileInfo.Name = fi.Name()
	fileInfo.Abs = filename
	fileInfo.Size = fi.Size()
	fileData := fi.Sys().(*syscall.Stat_t)
	// atime,ctime,Mtime分别是访问时间, 创建时间和修改时间
	fileInfo.AccessTimestamp = time.Unix(int64(fileData.Atim.Sec), int64(fileData.Atim.Nsec)).UnixNano() / 1e6
	fileInfo.CreateTimestamp = time.Unix(int64(fileData.Ctim.Sec), int64(fileData.Ctim.Nsec)).UnixNano() / 1e6
	fileInfo.ModifyTimestamp = time.Unix(int64(fileData.Mtim.Sec), int64(fileData.Mtim.Nsec)).UnixNano() / 1e6
	now := time.Now().UnixNano() / 1e6
	// 获取当前时间与访问时间/创建时间/修改时间比较, 一旦有一个时间在busyInterval以内则认为文件正在被传输中
	if AbsInt64(now-fileInfo.AccessTimestamp) < busyInterval || AbsInt64(now-fileInfo.CreateTimestamp) < busyInterval || AbsInt64(now-fileInfo.ModifyTimestamp) < busyInterval {
		return true, &fileInfo
	}
	return false, &fileInfo
}

// 判断文件是否被占用
func IsFileBusy(filename string) (bool, *vo.FileInfo) {
	return IsFileBusyByInterval(filename, s_const.BUSY_INTERVAL)
}
