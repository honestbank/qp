package handlers_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/qp/handlers"
)

func TestNewPrometheus(t *testing.T) {
	prometheus, err := handlers.PushToPrometheus("http://localhost:41129/", "qp")
	assert.NoError(t, err)
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(nil)
	prometheus(nil)
	prometheus(nil)
	prometheus(nil)
	prometheus(nil)
	prometheus(nil)
	prometheus(nil)
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	prometheus(errors.New("error"))
	time.Sleep(time.Second * 3)
}
