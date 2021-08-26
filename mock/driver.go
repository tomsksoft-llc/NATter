package mock

import (
	"context"

	"NATter/driver"
	"NATter/entity"

	nats "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/mock"
)

type DriverConn struct {
	mock.Mock
}

func (c *DriverConn) Serve(ctx context.Context) error {
	agrs := c.Called(ctx)

	return agrs.Error(0)
}

func (c *DriverConn) Close() error {
	agrs := c.Called()

	return agrs.Error(0)
}

func (c *DriverConn) Receiver(route *entity.Route) driver.Receiver {
	agrs := c.Called(route)

	return agrs.Get(0).(driver.Receiver)
}

func (c *DriverConn) Sender(route *entity.Route) driver.Sender {
	agrs := c.Called(route)

	return agrs.Get(0).(driver.Sender)
}

type DriverReceiver struct {
	mock.Mock
}

func (r *DriverReceiver) Listen(sender driver.Sender) error {
	args := r.Called(sender)

	return args.Error(0)
}

func (r *DriverReceiver) ListenRequest(sender driver.Sender) error {
	args := r.Called(sender)

	return args.Error(0)
}

type DriverSender struct {
	mock.Mock
}

func (s *DriverSender) Send(payload []byte) error {
	args := s.Called(payload)

	return args.Error(0)
}

func (s *DriverSender) Request(payload []byte) ([]byte, error) {
	args := s.Called(payload)

	return args.Get(0).([]byte), args.Error(1)
}

type DriverNatsConn struct {
	mock.Mock
}

func (c *DriverNatsConn) Subscribe(topic string, handler func(*nats.Msg) error) error {
	args := c.Called(topic, handler)

	return args.Error(0)
}

func (c *DriverNatsConn) Publish(topic string, payload []byte) error {
	args := c.Called(topic, payload)

	return args.Error(0)
}

func (c *DriverNatsConn) Request(topic string, payload []byte) (*nats.Msg, error) {
	args := c.Called(topic, payload)

	return args.Get(0).(*nats.Msg), args.Error(1)
}

type DriverKafkaConn struct {
	mock.Mock
}

func (c *DriverKafkaConn) Subscribe(topic string, handler func([]byte) error) error {
	args := c.Called(topic, handler)

	return args.Error(0)
}

func (c *DriverKafkaConn) Publish(topic string, payload []byte) error {
	args := c.Called(topic, payload)

	return args.Error(0)
}
