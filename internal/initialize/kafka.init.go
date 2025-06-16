package initialize

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"maps"

	"github.com/segmentio/kafka-go"
	"go-ecommerce-backend-api.com/global"
)

//initKafka producer

// var KafkaProducer *kafka.Writer

func InitKafka() {
	k := global.Config.Kafka
	brokerAddress := fmt.Sprintf("%s:%d", k.Host, k.Port)

	// Init map
	global.KafkaProducers = make(map[string]*kafka.Writer)

	// Kết nối controller để tạo topic
	conn, err := kafka.Dial("tcp", brokerAddress)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("❌ Failed to get Kafka controller: %v", err)
	}
	controllerAddr := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	controllerConn, err := kafka.Dial("tcp", controllerAddr)
	if err != nil {
		log.Fatalf("❌ Failed to dial controller: %v", err)
	}
	defer controllerConn.Close()

	// Duyệt qua tất cả topic trong config
	// Manually iterate over struct fields since k.Topic is a struct, not a map
	topics := make(map[string]string)
	maps.Copy(topics, k.Topic)

	for key, topicName := range topics {
		// Tự tạo topic nếu chưa có
		err := controllerConn.CreateTopics(kafka.TopicConfig{
			Topic:             topicName,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			log.Fatalf("❌ Failed to create topic %s: %v", topicName, err)
		}

		// Tạo writer và lưu vào map
		global.KafkaProducers[key] = &kafka.Writer{
			Addr:     kafka.TCP(brokerAddress),
			Topic:    topicName,
			Balancer: &kafka.LeastBytes{},
		}

		log.Printf("✅ Kafka producer ready for [%s] -> topic [%s]", key, topicName)
	}
}
func CloseKafka() {
	for key, writer := range global.KafkaProducers {
		if err := writer.Close(); err != nil {
			log.Printf("⚠️ Failed to close Kafka producer [%s]: %v", key, err)
		} else {
			log.Printf("✅ Closed Kafka producer [%s]", key)
		}
	}
}
