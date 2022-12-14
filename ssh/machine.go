package ssh

import (
	"fmt"
	sftp2 "github.com/QianLongGit/mskj.go/sftp"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/sftp"
	"strings"
	"time"
)

type RemoteMachine struct {
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// 获取sftp连接
func (m *RemoteMachine) GetSftpClient() (*sftp.Client, error) {
	return sftp2.GetSftpConnectClient(m.Username, m.Password, m.Ip, m.Port)
}

// 获取远端机器home目录
func (m *RemoteMachine) GetHomeDir() (string, error) {
	r := RemoteSshCommand(*m, false, []string{"pwd"})
	if r.Err != nil {
		return "", r.Err
	}

	homeDir := strings.TrimSpace(r.Result)
	return homeDir, nil
}

// 检查ssh工具是否安装
func (m *RemoteMachine) IsInstalledSsh() (bool, error) {
	r := RemoteSshCommand(*m, true, []string{
		"pwd",
	})
	if r.Err != nil {
		return false, r.Err
	}
	return true, nil
}

// 检查expect工具是否安装
func (m *RemoteMachine) IsInstalled(programName string) (bool, error) {
	r := RemoteSshCommand(*m, true, []string{
		"whereis " + programName,
	})
	if r.Err != nil {
		return false, r.Err
	}
	res := strings.Split(r.Result, ":")
	if len(res) == 2 {
		if strings.TrimSpace(res[1]) != "" {
			return true, nil
		}
	}
	return false, fmt.Errorf("未找到%s程序目录", programName)
}

// 重启机器
func (m *RemoteMachine) Reboot() {
	logs.Info(m.Ip, "开始重启操作系统")
	r := RemoteSshCommand(*m, false, []string{
		m.ExceptSudoWrapper("sudo reboot"),
	})
	for {
		time.Sleep(time.Second * 5)
		r = RemoteSshCommandAutoSessionCloseTimeout(*m, false, []string{"pwd"}, 2)
		if !r.IsConnect {
			logs.Info(m.Ip, "等待重连中...")
		} else {
			break
		}
	}
	logs.Info(m.Ip, "操作系统重启成功")
}

// 指定机器中的某个文件, 添加一行或多行字符
func (m *RemoteMachine) AddLinesToFile(filename string, lines ...string) error {
	var r Result
	for _, s := range lines {
		r = RemoteSshCommand(*m, true, []string{
			`cat ` + filename,
		})
		if r.Err != nil {
			continue
		}
		if strings.Contains(r.Result, s) {
			logs.Warn(fmt.Sprintf("%s上%s中已存在字符串[%s], 添加失败", m.Ip, filename, s))
		} else {
			r = RemoteSshCommand(*m, false, []string{
				fmt.Sprintf(`echo '%s' > /tmp/AddLinesToFile.txt`, s),
			})
			r = RemoteSshCommand(*m, false, []string{
				m.ExceptSudoWrapper("cat /tmp/AddLinesToFile.txt >> " + filename),
			})
			r = RemoteSshCommand(*m, false, []string{
				"rm -f /tmp/AddLinesToFile.txt",
			})
			if r.Err != nil {
				logs.Warn(fmt.Sprintf("为%s添加指定字符串[%s]到%s失败 %s", m.Ip, s, filename, r.Err))
			}
		}
	}
	if r.Err != nil {
		return r.Err
	}
	logs.Info(fmt.Sprintf("为%s添加指定字符串\n%s\n到%s成功", m.Ip, strings.Join(lines, "\n"), filename))
	return nil
}

// 使用expect命令包装sudo命令
func (m *RemoteMachine) ExceptSudoWrapper(cmd string) string {
	return m.ExceptSudoWrapperByExcept(cmd, "")
}

// 使用expect命令包装sudo命令
func (m *RemoteMachine) ExceptSudoWrapperByExcept(cmd string, expects ...string) string {
	expect := strings.TrimSpace(strings.Join(expects, " ; "))
	if expect == "" {
		expect = ""
	} else {
		// 个数不为空, 末尾加分号
		expect += " ; "
	}
	return `expect -c "spawn sudo sh -c \" ` + cmd + ` \" ; expect \"password for ` + m.Username + `:\" { send \"` + m.Password + `\r\" } ; ` + expect + ` set timeout -1 ; expect eof ;"`
}
