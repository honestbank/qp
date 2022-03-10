package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/honestbank/qp/queue/kafka"
	"time"
)

//func main() {
//	q, err := queue.NewKafka([]string{"localhost:9092"}, "topic-a", nil, "topic-a", "job-name", 20)
//	if err != nil {
//		panic(err)
//	}
//	job := qp.NewJob(q)
//	job.SetWorker(func(ev queue.Message) error {
//		msg := &sarama.ConsumerMessage{}
//		_ = ev.As(msg)
//		println(string(msg.Value))
//
//		return nil
//	})
//	job.Start()
//}


func producer1() {
	producer, _ := kafka.GetProducer([]string{"localhost:9092"}, "my-app")
	for i := 0; i < 100; i++ {
		producer.SendMessage(&sarama.ProducerMessage{
			Topic: "topic-a",
			Key:   nil,
			Value: sarama.StringEncoder(time.Now().Format(time.RFC3339)),
		})
		time.Sleep(time.Second)
	}
}

func main() {
	go func() {
		producer1()
	}()

	consumerGroup, err := kafka.GetConsumer([]string{"localhost:9092"}, "my-app", "my-group")
	if err != nil {
		panic(err)
	}

	defer func() { _ = consumerGroup.Close() }()

	// Track errors
	go func() {
		for err := range consumerGroup.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	channel := make(chan *kafka.KafkaMessage)
	go func(channel chan *kafka.KafkaMessage) {
		consumerGroup.Consume(context.Background(), []string{"topic-a"}, kafka.NewConsumer(channel))
	}(channel)

	for {
		fmt.Println((<-channel).String())
	}
}
