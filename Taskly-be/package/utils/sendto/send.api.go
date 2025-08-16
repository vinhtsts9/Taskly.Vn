package sendto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MailRequest struct {
	ToEmail     string `json:"toEmail"`
	MessageBody string `json:"messageBody"`
	Subject     string `json:"subject"`
	Attachment  string `json:"attachment"`
}

func SendEmailToJavaByApi(otp string, email string, purpose string) error {
	postURL := "http://localhost:8088/email/send_text"

	mailRequest := MailRequest{
		ToEmail:     email,
		MessageBody: "OTP is" + otp,
		Subject:     "Verify OTP" + purpose,
		Attachment:  "path/to/email",
	}

	requestBody, err := json.Marshal(mailRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content_Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Sprintln("Response status: ", resp.Status)
	return nil
}
