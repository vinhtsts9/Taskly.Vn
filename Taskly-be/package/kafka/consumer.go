package kafka

import (
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

type Consumer struct {
	client sarama.Consumer
}

func NewConsumer(brokers []string) (*Consumer, error) {
	client, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Kafka consumer initialized")
	return &Consumer{client: client}, nil
}

func (c *Consumer) Consume(topic string, handler func(message *sarama.ConsumerMessage) error) error {
	partitions, err := c.client.Partitions(topic)
	if err != nil {
		return err
	}
	for _, partition := range partitions {
		partitionConsumer, err := c.client.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}
		defer partitionConsumer.Close()
		// xu lys tin hieu ket thuc chuong trinh
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		for {
			select {
			case msg := <-partitionConsumer.Messages():
				err := handler(msg)
				if err != nil {
					log.Print("error handling message, %v", err)
				}
			case <-signals:
				log.Print("Interrupt signal received, shutting down consumer...")
				return nil
			}
		}

	}
	return nil
}

func (c *Consumer) Close() error {
	return c.client.Close()
}
