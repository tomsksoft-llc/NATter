package nats

import (
	"NATter/driver"
	"NATter/driver/msgbroker"

	nats "github.com/nats-io/nats.go"
)

type receiver struct {
	conn  Conn
	topic string
}

func (r *receiver) Listen(sender driver.Sender) error {
	return r.conn.Subscribe(r.topic, func(msg *nats.Msg) error {
		msgbroker.LogDebugReceived(r.topic, msg.Data)

		return sender.Send(msg.Data)
	})
}

func (r *receiver) ListenRequest(sender driver.Sender) error {
	return r.conn.Subscribe(r.topic, func(msg *nats.Msg) error {
		msgbroker.LogDebugReceived(r.topic, msg.Data)

		respb, err := sender.Request(msg.Data)

		if err != nil {
			return err
		}

		if err := msg.Respond(respb); err != nil {
			return msgbroker.ErrRespond(err, msg.Subject)
		}

		msgbroker.LogDebugResponded(msg.Subject, msg.Data)

		return nil
	})
}
