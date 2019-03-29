package service

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"github.com/toolkits/file"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Nssm struct {
	// nssm文件路径
	NssmPath string
	// 需要创建的服务名称
	ServiceName string
	// app程序文件路径
	AppFile string
	// app运行时工作路径
	AppDirectoryPath string
	// 状态0(未安装),1(正常运行),2(已暂停),-1(已停止)
	Status int
}

// 获取nssm实例
// nssmPath string nssm文件路径
// serviceName string 需要创建的服务名称
// appFile string app文件路径
// appDirectoryPath string app工作路径
func GetNewNssm(nssmPath string, serviceName string, appFile string, appDirectoryPath string) (*Nssm, error) {
	// nssmPath路径是否存在
	if !file.IsExist(nssmPath) {
		return nil, errors.New(fmt.Sprintf(`The nssm path "%s" is not exist`, nssmPath))
	}
	// nssmPath文件是否存在
	if !file.IsFile(nssmPath) {
		arch := 32
		if runtime.GOARCH == "amd64" {
			arch = 64
		}
		nssmPath = filepath.Join(nssmPath, fmt.Sprintf("nssm_%d.exe", arch))
		return GetNewNssm(nssmPath, serviceName, appFile, appDirectoryPath)
	}
	// 工作目录默认为当前目录
	root, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// appDirectoryPath目录是否存在
	if !file.IsExist(appDirectoryPath) {
		appDirectoryPath = filepath.Join(root, appDirectoryPath)
		if file.IsExist(appDirectoryPath) {
			return GetNewNssm(nssmPath, serviceName, appFile, appDirectoryPath)
		} else {
			return nil, errors.New("appDirectoryPath目录不存在")
		}
	}
	return &Nssm{
		NssmPath:         nssmPath,
		ServiceName:      serviceName,
		AppFile:          appFile,
		AppDirectoryPath: appDirectoryPath,
	}, nil
}

// 刷新状态
func (n *Nssm) RefreshStatus() error {
	// 执行命令
	cmd := fmt.Sprintf("%s status %s", n.NssmPath, n.ServiceName)
	logs.Debug(cmd)
	c := exec.Command("cmd", "/C", cmd)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	output := strings.ToLower(strings.TrimSpace(out.String()))
	if err != nil && output == "" {
		return err
	}
	logs.Error(output)
	if strings.Contains(output, "service_stopped") {
		// 已停止
		n.Status = -1
	} else if strings.Contains(output, "service_running") {
		// 正在运行
		n.Status = 1
	} else if strings.Contains(output, "service_paused") {
		// 已暂停
		n.Status = 2
	} else {
		// 未安装
		n.Status = 0
	}
	return nil
}

// 启动
func (n *Nssm) Start() error {
	n.RefreshStatus()
	if n.Status == 1 {
		return errors.New("服务正在运行, 无须重复启动")
	}
	if n.Status != -1 {
		n.Stop()
	}
	// 执行命令
	cmd := fmt.Sprintf("%s start %s", n.NssmPath, n.ServiceName)
	logs.Debug(cmd)
	c := exec.Command("cmd", "/C", cmd)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	output := strings.TrimSpace(out.String())
	if err != nil && output == "" {
		return err
	}
	n.Status = 1
	return nil
}

// 关闭
func (n *Nssm) Stop() error {
	n.RefreshStatus()
	if n.Status == -1 {
		return errors.New("服务已停止, 无须重复关闭")
	}
	// 执行命令
	cmd := fmt.Sprintf("%s stop %s", n.NssmPath, n.ServiceName)
	logs.Debug(cmd)
	c := exec.Command("cmd", "/C", cmd)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	output := strings.TrimSpace(out.String())
	if err != nil && output == "" {
		return err
	}
	n.Status = -1
	return nil
}

// 安装
func (n *Nssm) Install() error {
	n.RefreshStatus()
	if n.Status != 0 {
		return errors.New("服务已存在, 无须再次安装")
	}
	// 工作目录默认为当前目录
	root, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	if file.IsFile(filepath.Join(n.AppDirectoryPath, n.AppFile)) {
		// 优先使用工作目录下的APP
		n.AppFile = filepath.Join(n.AppDirectoryPath, n.AppFile)
	} else if file.IsFile(filepath.Join(root, n.AppFile)) {
		// 其次是运行目录
		n.AppFile = filepath.Join(root, n.AppFile)
	} else if !file.IsFile(n.AppFile) {
		return errors.New("appFile文件不存在")
	}
	// 执行命令
	cmd := fmt.Sprintf("%s install %s %s", n.NssmPath, n.ServiceName, n.AppFile)
	logs.Debug(cmd)
	c := exec.Command("cmd", "/C", cmd)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	output := strings.TrimSpace(out.String())
	if err != nil && output == "" {
		return err
	}
	if !strings.Contains(output, "installed successfully") {
		// 设置工作路径
		n.SetAppDirectory(n.AppDirectoryPath)
		// 安装成功, 默认没有运行
		n.Status = -1
	}
	return nil
}

// 卸载
func (n *Nssm) Remove() error {
	n.RefreshStatus()
	if n.Status == 0 {
		return errors.New("服务不存在, 无须卸载")
	}
	if n.Status != -1 {
		n.Stop()
	}
	// 执行命令
	cmd := fmt.Sprintf("%s remove %s confirm", n.NssmPath, n.ServiceName)
	logs.Debug(cmd)
	c := exec.Command("cmd", "/C", cmd)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	output := strings.TrimSpace(out.String())
	if err != nil && output == "" {
		return err
	}
	if strings.Contains(output, "removed successfully") {
		// 卸载成功
		n.Status = 0
	}
	return nil
}

// 设置工作目录
func (n *Nssm) SetAppDirectory(appDirectory string) error {
	n.RefreshStatus()
	if n.Status != 0 {
		return errors.New("服务不存在, 无法设置工作目录")
	}
	// appDirectory是否存在
	if !file.IsExist(appDirectory) || file.IsFile(appDirectory) {
		return errors.New(fmt.Sprintf(`The app path "%s" is not exist`, appDirectory))
	}
	// 执行命令
	cmd := fmt.Sprintf("%s set AppDirectory %s %s", n.NssmPath, n.ServiceName, appDirectory)
	c := exec.Command("cmd", "/C", cmd)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	output := strings.TrimSpace(out.String())
	if err != nil && output == "" {
		return err
	}
	n.AppDirectoryPath = appDirectory
	return nil
}
