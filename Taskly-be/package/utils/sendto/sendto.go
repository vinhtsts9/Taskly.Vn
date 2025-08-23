package sendto

import (
	"fmt"
	"net/smtp"
	"strings"

	"Taskly.com/m/global"

	"go.uber.org/zap"
)

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

func BuildMessage(email Mail) string {
	// build headers
	var b strings.Builder
	b.WriteString("MIME-Version: 1.0;\r\n")
	b.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n")
	b.WriteString(fmt.Sprintf("From: %s\r\n", email.From.Address))
	b.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ",")))
	b.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	b.WriteString("\r\n")
	b.WriteString(email.Body)
	return b.String()
}

// SendTextEmail gửi email dùng SMTP (Gmail App Password recommended).
// - to: danh sách email nhận
// - from: địa chỉ from (nên trùng SMTP_USERNAME để tránh bị block)
func SendTextEmail(to []string, from string, otp string) error {
	// lấy config từ env (không hardcode)
	host := global.ENVSetting.SMTP_HOST
	port := global.ENVSetting.SMTP_PORT
	username := global.ENVSetting.SMTP_USERNAME
	password := global.ENVSetting.SMTP_PASSWORD
	if host == "" || port == "" || username == "" || password == "" {
		err := fmt.Errorf("smtp config missing")
		global.Logger.Error("Email send failed::", zap.Error(err))
		return err
	}

	contentEmail := Mail{
		From:    EmailAddress{Address: from, Name: "Taskly"},
		To:      to,
		Subject: "OTP Verification",
		Body:    fmt.Sprintf("<p>Your otp is <b>%s</b>, please enter to verify your account</p>", otp),
	}

	messageEmail := BuildMessage(contentEmail)

	// PlainAuth: username is full email, host without port
	auth := smtp.PlainAuth("", username, password, host)

	addr := fmt.Sprintf("%s:%s", host, port)

	// dùng smtp.SendMail với addr (587) sẽ thực hiện STARTTLS nếu server hỗ trợ
	if err := smtp.SendMail(addr, auth, from, to, []byte(messageEmail)); err != nil {
		global.Logger.Error("Email send failed::", zap.Error(err))
		return err
	}

	return nil
}
