package qp

import "github.com/honestbank/qp/queue"

type Worker interface {
	SetWorker(func(ev queue.Message) error)
	OnResult(func(err error))
	Start()
}

type job struct {
	provider      queue.Queue
	job           func(ev queue.Message) error
	stopRequested bool
	onResult      func(err error)
}
