package queue

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type q struct {
	queueURL    *string
	queueClient *sqs.Client
}

type message struct {
	payload types.Message
	ack     func(handle *string) error
}

func (m message) String() string {
	return *m.payload.Body
}

func (m message) As(d interface{}) error {
	return json.Unmarshal([]byte(*m.payload.Body), &d)
}

func (m message) Ack() error {
	handle := m.payload.ReceiptHandle

	return m.ack(handle)
}

func (s q) Peek() (Message, error) {
	msg, err := s.queueClient.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:              s.queueURL,
		MaxNumberOfMessages:   1,
		WaitTimeSeconds:       5,
		VisibilityTimeout:     20,
		MessageAttributeNames: []string{"All"},
	})
	if err != nil {
		return nil, err
	}
	if len(msg.Messages) == 0 {
		return nil, nil
	}

	return message{
		payload: msg.Messages[0],
		ack: func(handle *string) error {
			_, err := s.queueClient.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{QueueUrl: s.queueURL, ReceiptHandle: handle})

			return err
		},
	}, nil
}

func NewSQSQueue(queueName string) (Queue, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := sqs.NewFromConfig(cfg)
	urlResult, err := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{QueueName: &queueName})
	if err != nil {
		return nil, err
	}

	return &q{
		queueURL:    urlResult.QueueUrl,
		queueClient: client,
	}, nil
}
