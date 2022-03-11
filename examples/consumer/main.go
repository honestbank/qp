package main

import (
	"context"
	"fmt"
	"github.com/honestbank/qp/examples"
	"github.com/honestbank/qp/queue/kafka"
	"time"
)

// func main() {
// 	q, err := queue.NewKafka([]string{"localhost:9092"}, topic, nil, "group_id", "job-name", 20)
// 	if err != nil {
// 		panic(err)
// 	}
// 	job := qp.NewJob(q)
// 	job.SetWorker(func(ev queue.Message) error {
// 		msg := &sarama.ConsumerMessage{}
// 		_ = ev.As(msg)
// 		println(string(msg.Value))
//
// 		return nil
// 	})
// 	job.Start()
// }

func main() {
	consumerGroup, err := kafka.GetConsumer([]string{"localhost:9092"}, "whatsapp", "group-b")
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

	msgChannel := make(chan *kafka.KafkaMessage)
	readyChannel := make(chan bool)
	go func(msgChannel chan *kafka.KafkaMessage, readyChannel chan bool) {
		for {
			err := consumerGroup.Consume(context.Background(), []string{examples.Topic}, kafka.NewConsumer(msgChannel, readyChannel))
			if err != nil {
				fmt.Println(err)
			}
		}
	}(msgChannel, readyChannel)

	for {
		msg := <-msgChannel
		fmt.Println((msg).String())
		time.Sleep(1*time.Second)
		msg.Ack()
		readyChannel <- true
	}
}
