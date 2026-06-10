// 我真诚地保证：
// 我自己独立地完成了整个程序从分析、设计到编码的所有工作。
// 如果在上述过程中，我遇到了什么困难而求教于人，那么，我将在程序实习报告中
// 详细地列举我所遇到的问题，以及别人给我的提示。
// 我的程序里中凡是引用到其他程序或文档之处，
// 例如教材、课堂笔记、网上的源代码以及其他参考书上的代码段,
// 我都已经在程序的注释里很清楚地注明了引用的出处。
// 我从未抄袭过别人的程序，也没有盗用别人的程序。
// 安俊豪
package email

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/smtp"

	"github.com/jordan-wright/email"
)

// Service 邮件服务
type Service struct {
	smtpHost string
	smtpPort int
	user     string
	password string
}

// New 创建邮件服务实例
func New() *Service {
	return &Service{
		smtpHost: "smtp.163.com",
		smtpPort: 465,
		user:     "18295659278@163.com",
		password: "RBxFtRLeR3WiNqa4",
	}
}

// GenerateCode 生成4位随机数字验证码
func (s *Service) GenerateCode() string {
	code := rand.Intn(9000) + 1000
	return fmt.Sprintf("%d", code)
}

// SendVerificationCode 发送验证码邮件
func (s *Service) SendVerificationCode(to, code, codeType string) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("图书管理系统 <%s>", s.user)
	e.To = []string{to}

	var subject, body string

	switch codeType {
	case "register":
		subject = "图书管理系统 - 用户注册验证码"
		body = fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #333;">图书管理系统 - 用户注册</h2>
				<p>您好！</p>
				<p>您正在进行用户注册，以下是您的验证码：</p>
				<div style="background: #f5f5f5; padding: 20px; text-align: center; font-size: 36px;
					  letter-spacing: 8px; font-weight: bold; color: #1890ff; border-radius: 4px;">
					%s
				</div>
				<p style="color: #999; font-size: 12px;">验证码有效期为5分钟，请尽快完成注册。</p>
				<p style="color: #999; font-size: 12px;">如非本人操作，请忽略此邮件。</p>
			</div>
		`, code)

	case "reset":
		subject = "图书管理系统 - 密码重置验证码"
		body = fmt.Sprintf(`
			<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto;">
				<h2 style="color: #333;">图书管理系统 - 密码重置</h2>
				<p>您好！</p>
				<p>您正在进行密码重置，以下是您的验证码：</p>
				<div style="background: #f5f5f5; padding: 20px; text-align: center; font-size: 36px;
					  letter-spacing: 8px; font-weight: bold; color: #1890ff; border-radius: 4px;">
					%s
				</div>
				<p style="color: #999; font-size: 12px;">验证码有效期为5分钟，请尽快完成重置。</p>
				<p style="color: #999; font-size: 12px;">如非本人操作，请忽略此邮件。</p>
			</div>
		`, code)
	}

	e.Subject = subject
	e.HTML = []byte(body)

	// 发送邮件（网易邮箱465端口需要SSL/TLS）
	addr := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
	auth := smtp.PlainAuth("", s.user, s.password, s.smtpHost)

	if err := e.SendWithTLS(addr, auth, &tls.Config{ServerName: s.smtpHost}); err != nil {
		return fmt.Errorf("邮件发送失败: %w", err)
	}
	return nil
}
