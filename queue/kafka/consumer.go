package kafka

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
)

type KafkaMessage struct {
	session *sarama.ConsumerGroupSession
	message *sarama.ConsumerMessage
}

type msg struct {
	Topic  string
	Header []*sarama.RecordHeader
	Body   string
	Key    string
}

func (k *KafkaMessage) String() string {
	m := msg{
		Topic:  k.message.Topic,
		Header: k.message.Headers,
		Body:   string(k.message.Value),
		Key:    string(k.message.Key),
	}
	b, _ := json.Marshal(m)

	return string(b)
}

func (k *KafkaMessage) As(d interface{}) error {
	val, ok := d.(*sarama.ConsumerMessage)
	if !ok {
		return errors.New("only `*sarama.ConsumerMessage` supported currently")
	}
	d = val

	return nil
}

func (k *KafkaMessage) Ack() error {
	(*k.session).MarkMessage(k.message, "")

	return nil
}

func (k *KafkaMessage) Nack() error {
	err := k.Ack()

	return err
	// todo: if that header as int > max
	// todo: produce into deadletter as is
	// todo: take the header as int, then increment and finally produce to our own topic
}

type consumer struct {
	messageChannel chan *KafkaMessage
	readyChannel   chan bool
}

func (c *consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Println(claim.InitialOffset())
		fmt.Println(msg.Offset)
		c.messageChannel <- &KafkaMessage{message: msg, session: &session}
		<- c.readyChannel
	}

	return nil
}

func NewConsumer(msgChannel chan *KafkaMessage, readyChannel chan bool) sarama.ConsumerGroupHandler {
	return &consumer{messageChannel: msgChannel, readyChannel: readyChannel}
}
