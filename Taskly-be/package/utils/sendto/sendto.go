package sendto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"Taskly.com/m/global"

	"go.uber.org/zap"
)

// Brevo API structs
type BrevoSender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type BrevoRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type BrevoPayload struct {
	Sender      BrevoSender      `json:"sender"`
	To          []BrevoRecipient `json:"to"`
	Subject     string           `json:"subject"`
	HtmlContent string           `json:"htmlContent"`
}

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

// SendTextEmail gửi email dùng Brevo API.
// - to: danh sách email nhận
// - from: địa chỉ from (phải là sender đã được xác thực trên Brevo)
func SendTextEmail(to []string, from string, otp string) error {
	// Lấy API key từ env
	apiKey := global.ENVSetting.Brevo_ApiKey
	if apiKey == "" {
		err := fmt.Errorf("brevo api key is missing (BREVO_API_KEY)")
		global.Logger.Error("Email send failed: Brevo API key not configured", zap.Error(err))
		return err
	}

	// Tạo danh sách người nhận cho Brevo
	var recipients []BrevoRecipient
	for _, emailAddr := range to {
		recipients = append(recipients, BrevoRecipient{Email: emailAddr})
	}

	// Tạo payload
	payload := BrevoPayload{
		Sender: BrevoSender{
			Name:  "Taskly", // Tên người gửi hiển thị
			Email: from,     // Email người gửi (phải được xác thực trên Brevo)
		},
		To:          recipients,
		Subject:     "Taskly - OTP Verification",
		HtmlContent: fmt.Sprintf("<html><body><p>Your OTP is <b>%s</b>. Please use it to verify your account.</p></body></html>", otp),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		global.Logger.Error("Failed to marshal Brevo payload", zap.Error(err))
		return err
	}

	// Tạo request
	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewBuffer(payloadBytes))
	if err != nil {
		global.Logger.Error("Failed to create Brevo request", zap.Error(err))
		return err
	}

	// Thêm headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	// Gửi request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		global.Logger.Error("Failed to send email via Brevo", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	// Kiểm tra response
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("brevo API returned non-2xx status: %d - %s", resp.StatusCode, string(body))
		global.Logger.Error("Brevo email send failed", zap.Error(err))
		return err
	}

	global.Logger.Info("Email sent successfully via Brevo", zap.Strings("to", to))
	return nil
}
