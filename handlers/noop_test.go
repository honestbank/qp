package handlers_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/honestbank/qp/handlers"
)

func TestNoop(t *testing.T) {
	t.Run("returns a function that's callable and doesn't panic", func(t *testing.T) {
		noop := handlers.Noop()
		assert.NotPanics(t, func() {
			noop(errors.New("text"))
		})
	})
}
