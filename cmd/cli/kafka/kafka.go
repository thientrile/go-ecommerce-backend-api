package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	kafka "github.com/segmentio/kafka-go"
)

const (
	kafkaURL   = "localhost:29092" // ⚠️ Dùng tên container trong Docker
	kafkaTopic = "user_topic_vip"
)

var kafkaProducer *kafka.Writer

// ===================== KAFKA CONFIG =====================

func getKafkaReader(brokerURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(brokerURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
	})
}

func getKafkaWriter(brokerURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// ===================== MODEL =====================

type StockInfo struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func newStock(message, msgType string) *StockInfo {
	return &StockInfo{
		Message: message,
		Type:    msgType,
	}
}

// ===================== HANDLERS =====================

func handleStockAction(c *gin.Context) {
	stock := newStock(c.Query("msg"), c.Query("type"))

	payload := map[string]interface{}{
		"action": "action",
		"info":   stock,
	}
	value, _ := json.Marshal(payload)

	msg := kafka.Message{
		Key:   []byte("action"),
		Value: value,
	}

	if err := kafkaProducer.WriteMessages(context.Background(), msg); err != nil {
		log.Printf("🔥 Kafka write error: %v", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "✅ Sent to Kafka successfully!"})
}

// ===================== CONSUMER =====================

func registerConsumerATC(id int) {
	groupID := fmt.Sprintf("consumer-group-%d", id)
	reader := getKafkaReader(kafkaURL, kafkaTopic, groupID)
	defer reader.Close()

	log.Printf("👂 Consumer %d is listening...\n", id)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("❌ Consumer %d error: %v", id, err)
			continue
		}
		fmt.Printf("🎯 Consumer %d → Topic: %s | Key: %s | Value: %s | Offset: %d | Partition: %d | Time: %s\n",
			id, m.Topic, string(m.Key), string(m.Value), m.Offset, m.Partition, m.Time.Format("15:04:05"))
	}
}

// ===================== MAIN =====================

func main() {
	kafkaProducer = getKafkaWriter(kafkaURL, kafkaTopic)
	defer kafkaProducer.Close()

	r := gin.Default()
	r.POST("/action/stock", handleStockAction)

	// Spawn multiple consumers
	go registerConsumerATC(1)
	go registerConsumerATC(2)
	go registerConsumerATC(3)
	go registerConsumerATC(4)

	log.Println("🚀 Server is running at http://localhost:8999")
	r.Run(":8999")
}
