package kafka

import (
	"github.com/Shopify/sarama"
)

func GetProducer(brokers []string, applicationName string) (sarama.SyncProducer, error) {
	producerConfig := kafkaProducerConfig(applicationName)

	return sarama.NewSyncProducer(brokers, producerConfig)
}

func GetConsumer(brokers []string, applicationName string, groupID string) (sarama.ConsumerGroup, error) {
	consumerConfig := kafkaConsumerConfig(applicationName)

	return sarama.NewConsumerGroup(brokers, groupID, consumerConfig)
}
