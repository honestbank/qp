package kafka

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
)

type Message struct {
	session      *sarama.ConsumerGroupSession
	message      *sarama.ConsumerMessage
	readyChannel chan bool
}

type msg struct {
	Topic  string
	Header []*sarama.RecordHeader
	Body   string
	Key    string
}

func (k *Message) String() string {
	m := msg{
		Topic:  k.message.Topic,
		Header: k.message.Headers,
		Body:   string(k.message.Value),
		Key:    string(k.message.Key),
	}
	b, _ := json.Marshal(m)

	return string(b)
}

func (k *Message) As(d interface{}) error {
	val, ok := d.(*sarama.ConsumerMessage)
	if !ok {
		return errors.New("only `*sarama.ConsumerMessage` supported currently")
	}
	d = val

	return nil
}

func (k *Message) Ack() error {
	(*k.session).MarkMessage(k.message, "")
	k.readyChannel <- true

	return nil
}

func (k *Message) Nack() error {
	err := k.Ack()

	return err
	// todo: if that header as int > max
	// todo: produce into deadletter as is
	// todo: take the header as int, then increment and finally produce to our own topic
}

type consumer struct {
	messageChannel chan *Message
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
		c.messageChannel <- &Message{message: msg, session: &session, readyChannel: c.readyChannel}
		<-c.readyChannel
	}

	return nil
}

func NewConsumer(msgChannel chan *Message, readyChannel chan bool) sarama.ConsumerGroupHandler {
	return &consumer{messageChannel: msgChannel, readyChannel: readyChannel}
}
