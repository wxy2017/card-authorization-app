package utils

import (
	"card-authorization/config"
	"card-authorization/log"
	"crypto/tls"
	"net/smtp"
	"strconv"
	"strings"
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
	// 发件人信息
	from := config.SystemConfig.EmailConfig.AuthEmail
	password := config.SystemConfig.EmailConfig.AuthPassword // 注意：这里必须使用QQ邮箱的SMTP授权码，而非登录密码
	smtpHost := config.SystemConfig.EmailConfig.SMTPHost
	smtpPort := strconv.Itoa(config.SystemConfig.EmailConfig.SMTPPort) // 465端口需要SSL加密

	// 构建邮件内容（必须包含标准头部）
	// 头部与正文之间需用空行(\r\n\r\n)分隔
	msg := []byte("To: " + to + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" + // 核心修复：指定为HTML格式
		"\r\n" + // 头部结束标记（必须有）
		body + "\r\n")

	// 构建认证信息（服务器地址不含端口）
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 465端口需要手动建立TLS连接
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	// 连接SMTP服务器（带SSL加密）
	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsConfig)
	if err != nil {
		log.Error("连接服务器失败：%v", err)
		return err
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		log.Error("创建客户端失败：%v", err)
		return err
	}
	defer client.Quit()

	// 进行认证
	if err := client.Auth(auth); err != nil {
		log.Error("认证失败：%v", err)
		return err
	}

	// 设置发件人
	if err := client.Mail(from); err != nil {
		log.Error("设置发件人失败：%v", err)
		return err
	}

	// 设置收件人（支持多个收件人）
	for _, addr := range strings.Split(to, ",") {
		if err := client.Rcpt(addr); err != nil {
			log.Error("设置收件人失败：%v", err)
			return err
		}
	}

	// 发送邮件内容
	data, err := client.Data()
	if err != nil {
		log.Error("准备发送内容失败：%v", err)
		return err
	}
	defer data.Close()

	// 写入邮件内容
	_, err = data.Write(msg)
	if err != nil {
		log.Error("发送内容失败：%v", err)
		return err
	}

	log.Info("邮件发送成功")
	return nil
}
