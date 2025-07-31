// package kafka/producer.go
package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type Producer struct {
	client sarama.SyncProducer
}

var (
	producerInstance *Producer
	once             sync.Once
)

// GetProducer trả về instance của Kafka producer (singleton)
func GetProducer(brokers []string) (*Producer, error) {
	var err error
	once.Do(func() {
		producerInstance, err = initKafkaProducer(brokers)
	})
	return producerInstance, err
}
func initKafkaProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	client, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	log.Println("Kafka producer initialized")
	return &Producer{client: client}, nil
}

// Send gửi tin nhắn tới Kafka
func (p *Producer) Send(topic string, message interface{}, maxRetries int) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	for i := 0; i < maxRetries; i++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(data),
		}

		_, _, err = p.client.SendMessage(msg)
		if err == nil {
			return nil
		}
		log.Printf("Retrying to send message to Kafka (%d/%d): %v", i+1, maxRetries, err)
	}

	return err
}
