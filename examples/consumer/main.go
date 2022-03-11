package main

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/honestbank/qp"
	"github.com/honestbank/qp/examples"
	"github.com/honestbank/qp/queue"
	"time"
)

func main() {
	val := "dead-letter-topic"
	q, err := queue.NewKafka([]string{"localhost:9092"}, examples.Topic, &val, "group_id", "job-name", 20)
	if err != nil {
		panic(err)
	}
	job := qp.NewJob(q)
	job.SetWorker(func(ev queue.Message) error {
		msg := &sarama.ConsumerMessage{}
		_ = ev.As(msg)
		println(ev.String())
		time.Sleep(time.Second * 5)

		return errors.New("some error")
	})
	job.Start()
}
