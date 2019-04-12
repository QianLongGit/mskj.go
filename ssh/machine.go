package ssh

import (
	"fmt"
	sftp2 "gitee.com/piupuer/go/sftp"
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
		`expect -c "spawn sudo reboot ; expect \"password for ` + m.Username + `:\" { send \"` + m.Password + `\r\" } ; set timeout -1 ; expect eof ;"`,
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

// 添加ubuntu的开机启动命令/etc/rc.local
func (m *RemoteMachine) AddUbuntuPoweredUpCmd(cmds ...string) error {
	var r Result
	for _, s := range cmds {
		r = RemoteSshCommand(*m, false, []string{
			fmt.Sprintf(`cat /etc/rc.local | grep '%s'`, s),
		})
		if r.Result != "" {
			logs.Warn("配置项[", s, "]已存在，将跳过该配置项的添加")
		} else {
			r = RemoteSshCommand(*m, false, []string{
				fmt.Sprintf(`echo '%s' > /tmp/AddUbuntuPoweredUpCmd.txt`, s),
			})
			r = RemoteSshCommand(*m, false, []string{
				`expect -c "spawn sudo bash -c \"cat /tmp/AddUbuntuPoweredUpCmd.txt >> /etc/rc.local \" ; expect \"password for ` + m.Username + `:\" { send \"` + m.Password + `\r\" } ; set timeout -1 ; expect eof ;"`,
			})
			r = RemoteSshCommand(*m, false, []string{
				`expect -c "spawn sudo bash -c \"rm - f /tmp/AddUbuntuPoweredUpCmd.txt \" ; expect \"password for ` + m.Username + `:\" { send \"` + m.Password + `\r\" } ; set timeout -1 ; expect eof ;"`,
			})
			if r.Err != nil {
				logs.Warn("配置项[", s, "]添加失败", r.Err)
			}
		}
	}
	if r.Err != nil {
		return r.Err
	}
	logs.Info(fmt.Sprintf("为%s添加/etc/rc.local开机启动项成功, 配置如下: \n%s", m.Ip, strings.Join(cmds, "\n")))
	return nil
}

// 添加环境变量到 /etc/profile 中
func (m *RemoteMachine) AddProfilesToEtcProfile(cmds ...string) error {
	var r Result
	for _, s := range cmds {
		r = RemoteSshCommand(*m, false, []string{
			fmt.Sprintf(`cat /etc/profile | grep '%s'`, s),
		})
		if r.Result != "" {
			logs.Warn("配置项[", s, "]已存在，将跳过该配置项的添加")
		} else {
			r = RemoteSshCommand(*m, false, []string{
				fmt.Sprintf(`echo '%s' > /tmp/AddProfilesToEtcProfile.txt`, s),
			})
			r = RemoteSshCommand(*m, false, []string{
				`expect -c "spawn sudo bash -c \"cat /tmp/AddProfilesToEtcProfile.txt >> /etc/profile \" ; expect \"password for ` + m.Username + `:\" { send \"` + m.Password + `\r\" } ; set timeout -1 ; expect eof ;"`,
			})
			r = RemoteSshCommand(*m, false, []string{
				`expect -c "spawn sudo bash -c \"rm - f /tmp/AddProfilesToEtcProfile.txt \" ; expect \"password for ` + m.Username + `:\" { send \"` + m.Password + `\r\" } ; set timeout -1 ; expect eof ;"`,
			})
			if r.Err != nil {
				logs.Warn("配置项[", s, "]添加失败", r.Err)
			}
		}
	}
	if r.Err != nil {
		return r.Err
	}
	logs.Info(fmt.Sprintf("为%s添加/etc/profile配置项添加成功, 配置如下: \n%s", m.Ip, strings.Join(cmds, "\n")))
	return nil
}
