package sendto

import (
	"encoding/json"
	"time"

	"Taskly.com/m/global"

	"github.com/segmentio/kafka-go"
)

func SendEmailKafka(email []string, otp string) error {
	body := make(map[string]interface{})

	body["otp"] = otp
	body["email"] = email

	bodyRequest, _ := json.Marshal(body)

	message := kafka.Message{
		Key:   []byte("otp-auth"),
		Value: []byte(bodyRequest),
		Time:  time.Now(),
	}

	err := global.KafkaProducer.Send("send-email", message, 3)
	if err != nil {
		return err
	}
	return nil
}
