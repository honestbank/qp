package kafka

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"strconv"
)

type Message struct {
	session         *sarama.ConsumerGroupSession
	message         *sarama.ConsumerMessage
	producer        sarama.SyncProducer
	deadLetterTopic *string
	maxReceiveCount uint
	readyChannel    chan bool
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
	switch d.(type) {
	case *sarama.ConsumerMessage:
		d = k.message

		return nil
	case *sarama.ProducerMessage:
		temp, err := k.convertToProducerMessage(k.message)
		d = temp

		return err
	default:

		return errors.New("only `*sarama.ConsumerMessage` supported currently")
	}
}

func (k *Message) convertToProducerMessage(source *sarama.ConsumerMessage) (*sarama.ProducerMessage, error) {
	return &sarama.ProducerMessage{
		Topic:   source.Topic,
		Key:     sarama.StringEncoder(source.Key),
		Value:   sarama.StringEncoder(source.Value),
		Headers: copyHeaders(source.Headers),
	}, nil
}

func copyHeaders(input []*sarama.RecordHeader) []sarama.RecordHeader {
	var result []sarama.RecordHeader
	if input == nil {
		return []sarama.RecordHeader{
			{
				Key:   []byte("message-receive-count"),
				Value: []byte(strconv.Itoa(0)),
			},
		}
	}
	for _, header := range input {
		if header == nil {
			continue
		}
		result = append(result, *header)
	}
	return result
}

func (k *Message) Ack() error {
	(*k.session).MarkMessage(k.message, "")
	k.readyChannel <- true

	return nil
}

func (k *Message) Nack() error {
	err := k.Ack()
	if err != nil {
		return err
	}
	if k.getReceiveCount() >= k.maxReceiveCount {
		return k.markMessageAsDead()
	}

	return k.nackReproduce()
}

func (k *Message) nackReproduce() error {
	producerMessage := sarama.ProducerMessage{}
	err := k.As(&producerMessage)
	if err != nil {
		return err
	}
	producerMessage.Headers = setReceiveCountIntoHeaders(producerMessage.Headers, k.getReceiveCount()+1)
	_, _, err = k.producer.SendMessage(&producerMessage)

	return err
}

func setReceiveCountIntoHeaders(headers []sarama.RecordHeader, targetCount uint) []sarama.RecordHeader {
	var targetHeader []sarama.RecordHeader
	for _, header := range headers {
		if string(header.Key) != "message-receive-count" {
			targetHeader = append(targetHeader, header)
			continue
		}
		targetHeader = append(targetHeader, sarama.RecordHeader{
			Key:   []byte("message-receive-count"),
			Value: []byte(strconv.Itoa(int(targetCount))),
		})
	}

	return targetHeader
}

func (k *Message) markMessageAsDead() error {
	if k.deadLetterTopic == nil {
		return nil
	}
	producerMessage := sarama.ProducerMessage{}
	err := k.As(&producerMessage)
	_, _, err = k.producer.SendMessage(&producerMessage)

	return err
}

func (k *Message) getReceiveCount() uint {
	msg := sarama.ConsumerMessage{}
	err := k.As(&msg)
	if err != nil {
		return 0
	}
	for _, header := range msg.Headers {
		val := getReceiveCountFromHeader(header.Key, header.Value)
		if val != nil {
			return *val
		}
	}

	return 0
}

func getReceiveCountFromHeader(key, value []byte) *uint {
	if string(key) != "message-receive-count" {
		return nil
	}
	val, err := strconv.Atoi(string(value))
	if err != nil {
		return nil
	}
	uintVal := uint(val)

	return &uintVal
}

type consumer struct {
	producer        sarama.SyncProducer
	deadLetterTopic *string
	messageChannel  chan *Message
	readyChannel    chan bool
	maxReceiveCount uint
}

func (c *consumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		c.messageChannel <- &Message{
			message:         msg,
			session:         &session,
			deadLetterTopic: c.deadLetterTopic,
			readyChannel:    c.readyChannel,
			maxReceiveCount: c.maxReceiveCount,
			producer:        c.producer,
		}
		<-c.readyChannel
	}

	return nil
}

func NewConsumer(msgChannel chan *Message, readyChannel chan bool, maxReceiveCount uint, deadLetterTopic *string, producer sarama.SyncProducer) sarama.ConsumerGroupHandler {
	return &consumer{
		producer:        producer,
		deadLetterTopic: deadLetterTopic,
		messageChannel:  msgChannel,
		readyChannel:    readyChannel,
		maxReceiveCount: maxReceiveCount,
	}
}
