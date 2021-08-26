package nats

type sender struct {
	conn  Conn
	topic string
}

func (s *sender) Send(payload []byte) error {
	return s.conn.Publish(s.topic, payload)
}

func (s *sender) Request(payload []byte) ([]byte, error) {
	msg, err := s.conn.Request(s.topic, payload)

	if err != nil {
		return nil, err
	}

	return msg.Data, nil
}
