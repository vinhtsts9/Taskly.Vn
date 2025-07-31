package sendto

import (
	"fmt"
	"net/smtp"
	"strings"

	"Taskly.com/m/global"

	"go.uber.org/zap"
)

const (
	SMTPHost     = "sandbox.smtp.mailtrap.io"
	SMTPPort     = "465"
	SMTPUsername = "0f856f53638897"
	SMTPPassword = "69b7a0dc9abfbe"
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
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", email.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", email.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", email.Body)

	return msg
}

func SendTextEmail(to []string, from string, otp string) error {
	contentEmail := Mail{
		From:    EmailAddress{Address: from, Name: "test"},
		To:      to,
		Subject: "OTP Verification",
		Body:    fmt.Sprintf("Your otp is %s, please enter to verify your account", otp),
	}

	messageEmail := BuildMessage(contentEmail)

	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)

	err := smtp.SendMail(SMTPHost+":465", auth, from, to, []byte(messageEmail))
	if err != nil {
		global.Logger.Error("Email send failed::", zap.Error(err))
		return err
	}

	return err
}
