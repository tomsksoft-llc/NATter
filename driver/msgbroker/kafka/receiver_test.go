package kafka

import (
	"testing"

	m "NATter/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestReceiverListen(t *testing.T) {
	conn := &m.DriverKafkaConn{}

	conn.
		On("Subscribe", "topic", mock.AnythingOfType("func([]uint8) error")).
		Run(func(args mock.Arguments) {
			err := args.Get(1).(func([]byte) error)([]byte("some-data"))

			assert.Nil(t, err)
		}).
		Return(nil)

	sender := &m.DriverSender{}

	sender.
		On("Send", []byte("some-data")).
		Return(nil)

	receiver := &receiver{
		conn:  conn,
		topic: "topic",
	}

	err := receiver.Listen(sender)

	assert.Nil(t, err)
}

func TestReceiverListenRequest(t *testing.T) {
	receiver := &receiver{}

	err := receiver.ListenRequest(&m.DriverSender{})

	assert.Error(t, err)
}
