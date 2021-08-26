package kafka

import (
	"testing"

	m "NATter/mock"

	"github.com/stretchr/testify/assert"
)

func TestSenderSend(t *testing.T) {
	conn := &m.DriverKafkaConn{}

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
	sender := &sender{}

	respb, err := sender.Request([]byte("request-data"))

	assert.Error(t, err)
	assert.Nil(t, respb)
}
