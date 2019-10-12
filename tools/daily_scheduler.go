//             ,%%%%%%%%,
//           ,%%/\%%%%/\%%
//          ,%%%\c "" J/%%%
// %.       %%%%/ o  o \%%%
// `%%.     %%%%    _  |%%%
//  `%%     `%%%%(__Y__)%%'
//  //       ;%%%%`\-/%%%'
// ((       /  `%%%%%%%'
//  \\    .'          |
//   \\  /       \  | |
//    \\/攻城狮保佑) | |
//     \         /_ | |__
//     (___________)))))))                   `\/'
/*
 * 修订记录:
 * long.qian 2019-03-07 16:48 创建
 */

/**
 * @author long.qian
 */

package tools

import (
	"fmt"
	"time"
)

//每天执行的定时器
//hour：24小时制，每天的几点开始启动
//f：启动时执行的方法
//param：方法参数
func DailyScheduler(hour, minute, second int, f func(param ...interface{}), param ...interface{}) {
	go func() {
		for {
			now := time.Now()
			var next time.Time
			temp := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, now.Location())
			if temp.Unix()-5 < now.Unix() {
				next = temp.Add(time.Hour * 24)
			} else {
				next = temp
			}
			n := next.Sub(now)
			t := time.NewTimer(n)
			fmt.Println(n.String())
			<-t.C
			f(param...)
		}
	}()
}
