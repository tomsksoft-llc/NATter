package driver

import (
	"context"

	"NATter/entity"
)

type Conn interface {
	Serve(context.Context) error
	Close() error

	Receiver(*entity.Route) Receiver
	Sender(*entity.Route) Sender
}

type Receiver interface {
	Listen(sender Sender) error
	ListenRequest(sender Sender) error
}

type Sender interface {
	Send(payload []byte) error
	Request(payload []byte) ([]byte, error)
}
