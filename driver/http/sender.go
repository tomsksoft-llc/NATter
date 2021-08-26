package http

import (
	"NATter/log"
)

type sender struct {
	endpoint string
}

func (s *sender) Send(payload []byte) error {
	_, err := request(s.endpoint, payload)

	return err
}

func (s *sender) Request(payload []byte) ([]byte, error) {
	respb, err := request(s.endpoint, payload)

	if err != nil {
		return nil, err
	}

	log.Debugf("received response from endpoint: %s", s.endpoint)

	return respb, nil
}
