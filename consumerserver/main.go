package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "book_events"
	brokerAddress := "localhost:9092"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{brokerAddress},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
		MaxWait:   time.Second,
	})
	defer r.Close()

	fmt.Println("Kafka consumer started. Waiting for messages...")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Error while reading message: %v", err)
		}
		fmt.Println("Message received: %s", string(msg.Value))
	}
}
