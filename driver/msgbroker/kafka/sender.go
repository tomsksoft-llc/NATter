package kafka

import (
	"github.com/pkg/errors"
)

type sender struct {
	conn  Conn
	topic string
}

func (s *sender) Send(payload []byte) error {
	return s.conn.Publish(s.topic, payload)
}

func (s *sender) Request(payload []byte) ([]byte, error) {
	return nil, errors.New("request is not supported by kafka")
}
