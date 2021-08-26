package nats

import (
	"testing"

	m "NATter/mock"

	nats "github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSenderSend(t *testing.T) {
	conn := &m.DriverNatsConn{}

	conn.
		On("Publish", "topic", []byte("some-data")).
		Return(nil)

	sender := &sender{
		conn:  conn,
		topic: "topic",
	}

	err := sender.Send([]byte("some-data"))

	assert.Nil(t, err)
}

func TestSenderRequest(t *testing.T) {
	conn := &m.DriverNatsConn{}

	conn.
		On("Request", "topic", []byte("request-data")).
		Return(&nats.Msg{
			Data: []byte("response-data"),
		}, nil)

	sender := &sender{
		conn:  conn,
		topic: "topic",
	}

	respb, err := sender.Request([]byte("request-data"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("response-data"), respb)
}

func TestSenderRequestOnError(t *testing.T) {
	conn := &m.DriverNatsConn{}

	conn.
		On("Request", "topic", []byte("request-data")).
		Return((*nats.Msg)(nil), errors.New("error"))

	sender := &sender{
		conn:  conn,
		topic: "topic",
	}

	respb, err := sender.Request([]byte("request-data"))

	assert.Error(t, err)
	assert.Nil(t, respb)
}
