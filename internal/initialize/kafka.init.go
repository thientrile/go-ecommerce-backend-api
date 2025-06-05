package initialize

import (
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
	"go-ecommerce-backend-api.com/global"
)

//initKafka producer

// var KafkaProducer *kafka.Writer

func InitKafka() {
	k := global.Config.Kafka
	s := fmt.Sprintf("%s:%d", k.Host, k.Port)
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP(s), // Replace with your Kafka broker address
		Topic:    k.Topic.Auth, // Replace with your Kafka topic name
		Balancer: &kafka.LeastBytes{},
	}
}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		log.Fatalf("Failed to close Kafka producer: %v", err)
	}
}
