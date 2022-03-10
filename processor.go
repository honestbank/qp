package qp

import (
	"os"
	"os/signal"
	"time"

	backoffPolicy "github.com/honestbank/backoff-policy"
	"github.com/honestbank/qp/handlers"
	"github.com/honestbank/qp/queue"
)

func NewJob(provider queue.Queue) Worker {
	return &job{
		provider: provider,
		onResult: handlers.Noop(),
	}
}

func (j *job) OnResult(f func(err error)) {
	j.onResult = f
}

func (j *job) SetWorker(f func(ev queue.Message) error) {
	j.job = f
}

func (j *job) Start() {
	if j.provider == nil {
		panic("no provider specified")
	}
	j.registerKillSignal()
	backoff := backoffPolicy.NewExponentialBackoffPolicy(time.Millisecond*100, 10)
	for ok := true; ok; ok = !j.stopRequested {
		backoff.Execute(j.eventLoop)
	}
}

func (j *job) registerKillSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			j.stopRequested = true
		}
	}()
}

func (j *job) eventLoop(marker backoffPolicy.Marker) {
	msg, err := j.provider.Peek()
	if err != nil {
		marker.MarkFailure()
		j.onResult(err)

		return
	}
	if msg == nil {
		// no messages yet, it's not a failure or a success so no marking required.

		return
	}
	err = j.job(msg)
	if err != nil {
		marker.MarkFailure()
		j.onResult(err)
		_ = msg.Nack()

		return
	}
	err = msg.Ack()
	if err != nil {
		marker.MarkFailure()
		j.onResult(err)

		return
	}
	j.onResult(nil)
	marker.MarkSuccess()
}
