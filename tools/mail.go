package tools

import (
	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
	"github.com/toolkits/file"
	"strings"
)

type Email struct {
	// 登录名称
	Username string `json:"username"`
	// 登录密码
	Password string `json:"password"`
	// 邮件主机
	Host string `json:"host"`
	// 邮件端口
	Port int `json:"port"`
	// 发件人名称
	FromName string `json:"from_name"`
	// 发件人
	From string `json:"from"`
	// 邮件主题
	Subject string
	// 纯文本
	Text string
	// Html格式
	HTML string
	// 收件人地址(1+)
	To []string
	// 抄送地址(1+)
	Cc []string
	// 暗送地址(1+)
	Bcc []string
	// 附件
	Attachment string
}

// 发送邮件
func SendEmail(email Email) error {
	m := gomail.NewMessage()

	fromName := strings.TrimSpace(email.FromName)
	if fromName == "" {
		fromName = "发件人"
	}
	from := strings.TrimSpace(email.From)
	if from == "" {
		return errors.New("发件人不得为空")
	}
	m.SetAddressHeader("From", from, fromName)

	if len(email.To) == 0 {
		return errors.New("收件人不得为空")
	}
	var tos []string
	for _, to := range email.To {
		tos = append(tos, m.FormatAddress(to, "收件人"))
	}
	m.SetHeader("To", tos...)
	var ccs []string
	for _, cc := range email.To {
		ccs = append(ccs, m.FormatAddress(cc, "收件人"))
	}
	m.SetHeader("Cc", ccs...)
	var bccs []string
	for _, bcc := range email.To {
		bccs = append(bccs, m.FormatAddress(bcc, "收件人"))
	}
	m.SetHeader("Bcc", bccs...)

	subject := strings.TrimSpace(email.Subject)
	if subject == "" {
		subject = "未设置标题"
	}
	m.SetHeader("Subject", subject)

	if strings.TrimSpace(email.HTML) != "" {
		// html
		m.SetBody("text/html", email.HTML)
	} else {
		// 纯文本
		m.SetBody("text/plain", email.Text)
	}

	attachment := strings.TrimSpace(email.Attachment)
	if file.IsFile(attachment) {
		// 添加附件
		m.Attach(attachment)
	}

	host := strings.TrimSpace(email.Host)
	if host == "" {
		return errors.New("邮件服务器不得为空")
	}
	port := email.Port
	if port == 0 {
		port = 465
	}
	username := strings.TrimSpace(email.Username)
	if username == "" {
		return errors.New("发件人登录名不得为空")
	}
	password := strings.TrimSpace(email.Password)
	if password == "" {
		return errors.New("发件人登录密码不得为空")
	}

	// 发送邮件服务器、端口、发件人账号、发件人密码
	d := gomail.NewDialer(host, port, username, password)
	return d.DialAndSend(m)
}
