package msgbroker

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrBadReply(t *testing.T) {
	assert.EqualError(t, ErrBadReply(errors.New("error"), "hello"), "bad reply from topic hello: error")
}

func TestErrPublish(t *testing.T) {
	assert.EqualError(t, ErrPublish(errors.New("error"), "hello"), "unable publish to topic hello: error")
}

func TestErrSubscribe(t *testing.T) {
	assert.EqualError(t, ErrSubscribe(errors.New("error"), "hello"), "unable subscribe to topic hello: error")
}

func TestErrUnsubscribe(t *testing.T) {
	assert.EqualError(t, ErrUnsubscribe(errors.New("error"), "hello"), "unable unsubscribe from topic hello: error")
}

func TestErrRespond(t *testing.T) {
	assert.EqualError(t, ErrRespond(errors.New("error"), "hello"), "unable respond to topic hello: error")
}
