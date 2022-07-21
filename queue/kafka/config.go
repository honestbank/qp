package kafka

import (
	"time"

	"github.com/Shopify/sarama"
)

func kafkaConsumerConfig(applicationName string) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.ClientID = applicationName
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}

func kafkaProducerConfig(applicationName string) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.ClientID = applicationName
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = time.Millisecond * 100
	config.Producer.Flush.Bytes = 64 * 1024
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Retry.Max = 5

	return config
}
