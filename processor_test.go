package qp_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/qp"
	"github.com/honestbank/qp/queue"
)

type mockQueue struct {
	returnNil         bool
	peekingShouldFail bool
	ackingShouldFail  bool
}
type mockMessage struct {
	acked            bool
	ackingShouldFail bool
}

func (m *mockMessage) String() string {
	if m.acked {
		return "true"
	}
	return "false"
}

func (m *mockMessage) As(d interface{}) error {
	return nil
}

func (m *mockMessage) Ack() error {
	m.acked = true
	if m.ackingShouldFail {
		return errors.New("acking error")
	}

	return nil
}

func (m *mockQueue) Peek() (queue.Message, error) {
	m.returnNil = !m.returnNil
	if m.returnNil {
		return nil, nil
	}
	if m.peekingShouldFail {
		return nil, errors.New("peeking error")
	}

	return &mockMessage{acked: false, ackingShouldFail: m.ackingShouldFail}, nil
}

func TestQp(t *testing.T) {
	t.Run("panics if no queue is defined", func(t *testing.T) {
		j := qp.NewJob(nil)
		assert.Panics(t, func() {
			j.Start()
		})
	})
	t.Run("can define a worker & result processor which is then called", func(t *testing.T) {
		job := qp.NewJob(&mockQueue{})
		workerCalled := 0
		resultProcessorCalled := 0
		job.SetWorker(func(ev queue.Message) error {
			workerCalled += 1
			time.Sleep(time.Millisecond * 200)

			return nil
		})
		job.OnResult(func(err error) {
			resultProcessorCalled += 1
		})
		go job.Start()
		time.Sleep(time.Second * 2)
		assert.True(t, workerCalled > 0)
		assert.True(t, resultProcessorCalled > 0)
	})
	t.Run("fail case while processing", func(t *testing.T) {
		job := qp.NewJob(&mockQueue{})
		resultProcessorCalled := 0
		job.SetWorker(func(ev queue.Message) error {
			return errors.New("some error")
		})
		job.OnResult(func(err error) {
			if err != nil {
				resultProcessorCalled += 1
			}
		})
		go job.Start()
		time.Sleep(time.Second * 2)
		assert.True(t, resultProcessorCalled > 0)
	})
	t.Run("fail case while peeking", func(t *testing.T) {
		job := qp.NewJob(&mockQueue{peekingShouldFail: true})
		resultProcessorCalled := 0
		job.SetWorker(func(ev queue.Message) error {
			time.Sleep(time.Millisecond * 200)

			return nil
		})
		job.OnResult(func(err error) {
			if err != nil {
				resultProcessorCalled += 1
			}
		})
		go job.Start()
		time.Sleep(time.Second * 2)
		assert.True(t, resultProcessorCalled > 0)
	})
	t.Run("fail case while acknowledging", func(t *testing.T) {
		job := qp.NewJob(&mockQueue{ackingShouldFail: true})
		resultProcessorCalled := 0
		job.SetWorker(func(ev queue.Message) error {
			time.Sleep(time.Millisecond * 200)

			return nil
		})
		job.OnResult(func(err error) {
			if err != nil {
				resultProcessorCalled += 1
			}
		})
		go job.Start()
		time.Sleep(time.Second * 2)
		assert.True(t, resultProcessorCalled > 0)
	})
}
