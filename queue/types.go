package queue

import "fmt"

type Queue interface {
	Peek() (Message, error)
}

type Message interface {
	fmt.Stringer
	As(d interface{}) error
	Ack() error
	Nack() error
}
