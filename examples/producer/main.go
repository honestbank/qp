package main

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"

	"github.com/honestbank/qp/examples"
	"github.com/honestbank/qp/queue/kafka"
)

func main() {
	producer, _ := kafka.GetProducer([]string{"localhost:9092"}, "my-app")
	for i := 0; i < 1000; i++ {
		producer.SendMessage(&sarama.ProducerMessage{
			Topic: examples.Topic,
			Key:   nil,
			Value: sarama.StringEncoder(time.Now().Format(time.RFC3339)),
		})
		fmt.Println(time.Now().Format(time.RFC3339))
		time.Sleep(time.Millisecond * 500)
	}
}
