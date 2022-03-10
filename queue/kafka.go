package queue

import (
	"context"

	"github.com/Shopify/sarama"

	"github.com/honestbank/qp/queue/kafka"
)

func NewKafka(brokers []string, topic string, deadLetterTopic *string, groupID string, applicationName string, maxReceiveCount int) (Queue, error) {
	producer, err := kafka.GetProducer(brokers, applicationName)
	if err != nil {
		return nil, err
	}

	consumerGroup, err := kafka.GetConsumer(brokers, applicationName, groupID)
	if err != nil {
		return nil, err
	}
	messageChan := make(chan *kafka.KafkaMessage)
	err = consumerGroup.Consume(context.Background(), []string{topic}, kafka.NewConsumer(messageChan))
	if err != nil {
		return nil, err
	}

	return kafkaQueue{
		producer:        producer,
		consumer:        consumerGroup,
		deadLetterTopic: deadLetterTopic,
		groupID:         groupID,
		maxReceiveCount: maxReceiveCount,
		topic:           topic,
		messages:        messageChan,
	}, nil
}

type kafkaQueue struct {
	producer        sarama.SyncProducer
	consumer        sarama.ConsumerGroup
	messages        chan *kafka.KafkaMessage
	deadLetterTopic *string
	groupID         string
	maxReceiveCount int
	topic           string
}

func (k kafkaQueue) Peek() (Message, error) {
	msg := <-k.messages

	return msg, nil
}
