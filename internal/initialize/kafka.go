// initialize/kafka.go
package initialize

import (
	"log"

	"Taskly.com/m/global"
	"Taskly.com/m/package/kafka"
)

func InitKafka() {
	stringBroker := global.Config.KafkaBroker.Brokers
	global.Logger.Sugar().Info(stringBroker)
	brokers := []string{stringBroker}
	producer, err := kafka.GetProducer(brokers)
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)
	}
	global.KafkaProducer = producer

	global.KafkaConsumer, err = kafka.NewConsumer(brokers)
	if err != nil {
		log.Fatalf("Error initializing Kafka consumer: %v", err)
	}
}
