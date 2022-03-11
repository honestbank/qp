package main

import (
	"github.com/Shopify/sarama"
	"github.com/honestbank/qp"
	"github.com/honestbank/qp/examples"
	"github.com/honestbank/qp/queue"
	"time"
)

func main() {
	q, err := queue.NewKafka([]string{"localhost:9092"}, examples.Topic, nil, "group_id", "job-name", 20)
	if err != nil {
		panic(err)
	}
	job := qp.NewJob(q)
	job.SetWorker(func(ev queue.Message) error {
		msg := &sarama.ConsumerMessage{}
		_ = ev.As(msg)
		println(ev.String())
		println(msg.Offset)
		time.Sleep(time.Second * 2)

		return nil
	})
	job.Start()
}

//func main() {
//	consumerGroup, err := kafka.GetConsumer([]string{"localhost:9092"}, "whatsapp", "group-b")
//	if err != nil {
//		panic(err)
//	}
//
//	defer func() { _ = consumerGroup.Close() }()
//
//	// Track errors
//	go func() {
//		for err := range consumerGroup.Errors() {
//			fmt.Println("ERROR", err)
//		}
//	}()
//
//	msgChannel := make(chan *kafka.Message)
//	readyChannel := make(chan bool)
//	go func(msgChannel chan *kafka.Message, readyChannel chan bool) {
//		for {
//			err := consumerGroup.Consume(context.Background(), []string{examples.Topic}, kafka.NewConsumer(msgChannel, readyChannel))
//			if err != nil {
//				fmt.Println(err)
//			}
//		}
//	}(msgChannel, readyChannel)
//
//	for {
//		msg := <-msgChannel
//		fmt.Println((msg).String())
//		time.Sleep(5*time.Second)
//		msg.Ack()
//	}
//}
