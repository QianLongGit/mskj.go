package tools

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

// 捕获错误信息写入文件
func CatchErrors() {
	errs := recover()
	if errs == nil {
		return
	}
	// 获取程序名称
	abs := filepath.Dir(os.Args[0])
	programName := filepath.Base(os.Args[0])
	// 获取当前时间
	now := GetNowDateFormat()
	nowTimestamp := time.Now().Unix()
	// 获取进程ID
	pid := os.Getpid()

	// 保存错误信息文件名:程序名_进程ID_当前时间
	logFile := filepath.Join(abs, "logs", fmt.Sprintf("%s_%d_%s_%d_dump.log", programName, pid, now, nowTimestamp))
	logs.Error("程序非正常崩溃: ", errs)
	logs.Error(fmt.Sprintf("错误详情日志已写入%s文件", logFile))
	_, err := CreateFileIfNotExists(logFile)
	f, err := os.Create(logFile)
	if err != nil {
		return
	}
	defer f.Close()

	// 输出panic信息
	f.WriteString("panic:\r\n")
	f.WriteString(fmt.Sprintf("%v", errs))
	f.WriteString("\r\nstack:\r\n")
	// 输出堆栈信息
	f.WriteString(string(debug.Stack()))
}
