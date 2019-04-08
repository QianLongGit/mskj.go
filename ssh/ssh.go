package ssh

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
)

type Result struct {
	IsConnect bool   `json:"is_connect"`
	Result    string `json:"Result"`
	Err       error  `json:"err"`
}

// 判断命令是否是允许的安全命令
func IsSafetyCmd(cmd string) (string, bool) {
	if strings.Contains(cmd, "rm") {
		if len(strings.Split(cmd, "/")) <= 1 {
			return fmt.Sprintf("rm命令 %s 不能删除小于2级的文件", cmd), false
		}
	}
	return "", true
}

// 执行远程命令
// isPrintDetail: 是否控制台自动打印命令执行过程信息（注：mac平台，可能打印出来的都是空白字符）
func RemoteSshCommand(machine RemoteMachine, isPrintDetail bool, cmds []string) Result {
	return RemoteSshCommandAutoSessionCloseTimeout(machine, isPrintDetail, cmds, 0)
}

// 执行远程命令，超时时，自动关闭session，sessionCloseSecond=0 代表不设置超时强制关闭session
// isPrintDetail:是否控制台自动打印命令执行过程信息（注：mac平台，可能打印出来的都是空白字符）
func RemoteSshCommandAutoSessionCloseTimeout(machine RemoteMachine, isPrintDetail bool, cmds []string, sessionCloseSecond int64) Result {
	var session *ssh.Session
	client, err := GetSshClient(machine.Username, machine.Password, machine.Ip, machine.Port)
	if err != nil {
		return Result{
			IsConnect: false,
			Err:       err,
		}
	}
	// create session
	if session, err = client.NewSession(); err != nil {
		return Result{
			IsConnect: false,
			Err:       fmt.Errorf("Failed New SSH Session To %s : %s", machine.Ip, err),
		}
	}

	if err != nil {
		return Result{
			IsConnect: false,
			Err:       err,
		}
	}
	defer closeClient(session, client)

	go func() {
		if sessionCloseSecond > 0 {
			sleep, err := time.ParseDuration(fmt.Sprintf("%ds", sessionCloseSecond))
			if err != nil {
				logs.Error("SSH自动session关闭失败：%s", err)
				return
			}
			time.Sleep(sleep)
			closeClient(session, client)
		}
	}()

	command := ""

	for i, cmd := range cmds {
		if msg, ok := IsSafetyCmd(cmd); !ok {
			return Result{
				IsConnect: true,
				Err:       fmt.Errorf(msg),
			}
		}
		if i == 0 {
			command = cmd
		} else {
			command = command + " && " + cmd
		}
	}

	var e bytes.Buffer
	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = &e
	var result string
	if command != "" {
		if isPrintDetail {
			go func() {
				time.Sleep(time.Second)
				lastS := ""
				for {
					if b.String() == "Exec Done!" {
						break
					}
					r := b.String()
					if lastS == "" {
						logs.Info(machine.Ip, ":")
						fmt.Println(r)
						lastS = r
					} else if lastS != r {
						logs.Info(machine.Ip, ":")
						fmt.Println(r)
						lastS = r
					}
					time.Sleep(time.Second * 3)
				}
			}()
		}
		if err := session.Run(command); err != nil {
			return Result{
				IsConnect: true,
				Err:       fmt.Errorf("命令 %s 执行失败：%s\n%s", command, err, e.String()),
				Result:    e.String(),
			}
		} else {
			result = b.String()
			b.Reset()
			b.WriteString("Exec Done!")
		}
		logs.Debug(command)
	}
	return Result{
		Result:    result,
		Err:       nil,
		IsConnect: true,
	}
}

func closeClient(session *ssh.Session, client *ssh.Client) {
	err := client.Close()
	if err != nil {
		logs.Error("SSH Client Close Error", err)
	}
	// 必须关闭Client，才能释放该ssh连接句柄
	session.Close()
}

func GetSshClient(user, password, host string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 8 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, fmt.Errorf("Connect To %s Error : %s", addr, err)
	}
	return client, nil
}
