package kafka

import (
	"NATter/driver"
	"NATter/driver/msgbroker"

	"github.com/pkg/errors"
)

type receiver struct {
	conn  Conn
	topic string
}

func (r *receiver) Listen(sender driver.Sender) error {
	return r.conn.Subscribe(r.topic, func(payload []byte) error {
		msgbroker.LogDebugReceived(r.topic, payload)

		return sender.Send(payload)
	})
}

func (r *receiver) ListenRequest(driver.Sender) error {
	return errors.New("response is not supported by kafka")
}
