package main

import (
	"context"
	"fmt"
	"github.com/honestbank/qp/examples"
	"github.com/honestbank/qp/queue/kafka"
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
	consumerGroup, err := kafka.GetConsumer([]string{"localhost:9092"}, "my-app", "group-a")
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
		consumerGroup.Consume(context.Background(), []string{examples.Topic}, kafka.NewConsumer(channel))
	}(channel)

	for {
		fmt.Println((<-channel).String())
	}
}
