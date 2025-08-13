package utils

import (
	"card-authorization/config"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strconv"
)

// EmailConfig 邮件配置结构体
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	AuthEmail    string
	AuthPassword string
}

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	// 邮件内容
	fromEmail := config.SystemConfig.EmailConfig.AuthEmail
	fromSMTPHost := config.SystemConfig.EmailConfig.SMTPHost
	fromSMTPPort := config.SystemConfig.EmailConfig.SMTPPort
	fromAuthPassword := config.SystemConfig.EmailConfig.AuthPassword
	fromAddr := fromSMTPHost + strconv.Itoa(fromSMTPPort)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body + "\r\n")
	// SMTP 认证
	auth := smtp.PlainAuth("", fromEmail, fromAuthPassword, fromSMTPHost+":"+strconv.Itoa(fromSMTPPort))
	// 发送邮件
	err := smtp.SendMail(fromAddr, auth, fromEmail, []string{to}, msg)
	if err != nil {
		return err
	}
	return nil
}

// SendEmail 发送邮件
func SendEmail2(to, subject, body string) error {
	from := config.SystemConfig.EmailConfig.AuthEmail
	fromSMTPHost := config.SystemConfig.EmailConfig.SMTPHost
	fromSMTPPort := config.SystemConfig.EmailConfig.SMTPPort
	fromAuthPassword := config.SystemConfig.EmailConfig.AuthPassword
	fromAddr := fromSMTPHost + strconv.Itoa(fromSMTPPort)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
		"\r\n" +
		body + "\r\n")

	// SMTP 认证
	addr := fromAddr
	auth := smtp.PlainAuth("", from, fromAuthPassword, fromSMTPHost)

	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         fromSMTPHost,
	}

	// 建立 SMTP 连接
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("连接 SMTP 服务器失败: %w", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, fromSMTPHost)
	if err != nil {
		return fmt.Errorf("创建 SMTP 客户端失败: %w", err)
	}
	defer client.Close()

	// 认证
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 认证失败: %w", err)
	}

	// 设置发件人
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %w", err)
	}

	// 设置收件人
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %w", err)
	}

	// 发送邮件内容
	data, err := client.Data()
	if err != nil {
		return fmt.Errorf("打开数据流失败: %w", err)
	}
	defer data.Close()

	if _, err := data.Write(msg); err != nil {
		return fmt.Errorf("写入邮件内容失败: %w", err)
	}

	// 关闭连接
	if err := client.Quit(); err != nil {
		return fmt.Errorf("关闭 SMTP 连接失败: %w", err)
	}

	return nil
}
