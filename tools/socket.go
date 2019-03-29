package tools

import (
	"fmt"
	"net"
	"time"
)

//判断指定的IP端口能否连接
func IsTcpConnected(ip string, port int, timeout time.Duration) bool {
	server := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", server, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
