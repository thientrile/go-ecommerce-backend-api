package sendto

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"go-ecommerce-backend-api.com/global"
)

func SendMessageToKafka(key string, value string) error {

	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
	}

	err := global.KafkaProducer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return err
	}

	return nil
}
