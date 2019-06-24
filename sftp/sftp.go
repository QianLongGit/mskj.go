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
 * long.qian 2018-08-20 09:44 创建
 */

/**
 * @author long.qian
 */

package sftp

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/sftp"
	"github.com/toolkits/file"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"mskj.go/tools"
	"os"
	"path"
	"time"
)

func GetSftpConnectClient(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

// sftp上传文件。注：需要手动关闭连接
func UploadFile(sftpClient *sftp.Client, localFilePath string, remoteDir string) error {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)

	err = sftpClient.MkdirAll(remoteDir)
	if err != nil {
		return err
	}
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		logs.Error("create sftpClient error : %s", path.Join(remoteDir, remoteFileName))
		return err
	}
	defer dstFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	var send int64
	var lastProcess string
	var showLog bool
	var success bool
	go func() {
		for {
			time.Sleep(time.Second * 5)
			showLog = true
			if success {
				break
			}
		}
	}()
	for {
		buf := make([]byte, 10240)
		n, rErr := srcFile.Read(buf)
		if rErr == io.EOF {
			break
		}
		_, err = dstFile.Write(buf)
		if err != nil {
			return err
		}
		send += int64(n)
		process := tools.PercentageFormat(send, srcInfo.Size())
		if showLog && lastProcess != process {
			lastProcess = process
			showLog = false
			logs.Info(fmt.Sprintf("本地文件%s已发送%s, 总共%s, 发送进度%s", localFilePath, tools.FileSizeFormat(send), tools.FileSizeFormat(srcInfo.Size()), process))
		}
	}

	success = true
	logs.Info(fmt.Sprintf("本地文件%s发送到远端服务器%s成功", localFilePath, remoteDir))
	return nil
}

// 上传文件夹。注：需要手动关闭连接
func UploadDirectory(sftpClient *sftp.Client, localDir string, remoteDir string) error {
	localFiles, err := ioutil.ReadDir(localDir)
	if err != nil {
		return err
	}

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localDir, backupDir.Name())
		remoteFilePath := path.Join(remoteDir, backupDir.Name())
		if backupDir.IsDir() {
			err = UploadDirectory(sftpClient, localFilePath, remoteFilePath)
			if err != nil {
				return err
			}
		} else {
			err = UploadFile(sftpClient, path.Join(localDir, backupDir.Name()), remoteDir)
			if err != nil {
				return err
			}
		}
	}
	logs.Info(fmt.Sprintf("本地路径%s发送到远端服务器路径%s成功", localDir, remoteDir))
	return nil
}

// 下载远程目录中的所有文件，到指定的本地目录
func DownloadDirectory(sftpClient *sftp.Client, localDir string, remoteDir string) error {
	remoteFiles, err := sftpClient.ReadDir(remoteDir)
	if err != nil {
		return err
	}

	for _, backupDir := range remoteFiles {
		localFilePath := path.Join(localDir, backupDir.Name())
		remoteFilePath := path.Join(remoteDir, backupDir.Name())
		if backupDir.IsDir() {
			err = DownloadDirectory(sftpClient, localFilePath, remoteFilePath)
			if err != nil {
				return err
			}
		} else {
			err = DownloadFile(sftpClient, localFilePath, remoteFilePath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// 下载文件，若本地文件夹不存在，自动创建
func DownloadFile(sftpClient *sftp.Client, localFilePath string, remoteDir string) error {
	srcFile, err := sftpClient.Open(remoteDir)
	if err != nil {
		logs.Error(err)
		return err
	}
	defer srcFile.Close()

	localDir := file.Dir(localFilePath)
	if !file.IsExist(localDir) {
		if err := os.MkdirAll(localDir, os.FileMode(0755)); err != nil {
			return fmt.Errorf("目录创建失败: %s", err)
		}
	}
	var localFileName = file.Basename(localFilePath)
	dstFile, err := os.Create(path.Join(localDir, localFileName))
	if err != nil {
		logs.Error(err)
		return err
	}
	defer dstFile.Close()

	if _, err = srcFile.WriteTo(dstFile); err != nil {
		logs.Error(err)
		return err
	}
	logs.Info(fmt.Sprintf("远端服务器路径%s下载到本地路径%s成功", remoteDir, localDir))
	return nil
}
