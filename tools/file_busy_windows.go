package tools

import (
	"fmt"
	"mskj.go/s_const"
	"mskj.go/vo"
	"os"
	"syscall"
	"time"
)

// 运行时: Windows

// 判断文件是否被占用
func IsFileBusyByInterval(filename string, busyInterval int64) (bool, *vo.FileInfo) {
	// 只能在1毫秒到1小时之间
	if busyInterval < 0 || busyInterval > 360000 {
		busyInterval = s_const.BUSY_INTERVAL
	}
	f, err := os.Open(filename)
	defer func() {
		if f != nil {
			err := f.Close()
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
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
	fileData := fi.Sys().(*syscall.Win32FileAttributeData)
	// LastAccessTime,CreationTime,LastWriteTime分别是访问时间, 创建时间和修改时间
	fileInfo.AccessTimestamp = fileData.LastAccessTime.Nanoseconds() / 1e6
	fileInfo.CreateTimestamp = fileData.CreationTime.Nanoseconds() / 1e6
	fileInfo.ModifyTimestamp = fileData.LastWriteTime.Nanoseconds() / 1e6
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
