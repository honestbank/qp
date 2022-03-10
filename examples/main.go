package main

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"

	"github.com/honestbank/qp"
	"github.com/honestbank/qp/queue"
	"github.com/honestbank/qp/queue/kafka"
)

func main() {
	q, err := queue.NewKafka([]string{"localhost:9092"}, "application-form-phone-saved", nil, "application-form-phone-saved", "job-name", 20)
	if err != nil {
		panic(err)
	}
	job := qp.NewJob(q)
	job.SetWorker(func(ev queue.Message) error {
		msg := &sarama.ConsumerMessage{}
		_ = ev.As(msg)
		println(string(msg.Value))

		return nil
	})
	job.Start()
}

// func temp() {
// 	producer, _ := kafka.GetProducer([]string{"localhost:9092"}, "my-app")
// 	for i := 0; i < 100; i++ {
// 		producer.SendMessage(&sarama.ProducerMessage{
// 			Topic: "application-form-phone-saved",
// 			Key:   nil,
// 			Value: sarama.StringEncoder(time.Now().Format(time.RFC3339)),
// 		})
// 		time.Sleep(time.Second)
// 	}
// }
func temp2() {
	producer, _ := kafka.GetConsumer([]string{"localhost:9092"}, "my-app", "1")
	channel := make(chan *kafka.KafkaMessage)
	go func(channel chan *kafka.KafkaMessage) {
		producer.Consume(context.Background(), []string{"application-form-phone-saved"}, kafka.NewConsumer(channel))
	}(channel)
	fmt.Println((<-channel).String())
}
