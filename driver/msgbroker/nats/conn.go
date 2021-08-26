package nats

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"NATter/driver"
	"NATter/driver/msgbroker"
	"NATter/entity"
	"NATter/errtpl"
	"NATter/log"

	nats "github.com/nats-io/nats.go"
)

const (
	DriverName = "nats"

	requestTimeout = time.Second * 10
)

type ConnConfig struct {
	Servers []string
	Token   string
	Group   string
	Name    string
}

type Conn interface {
	Subscribe(topic string, handler func(*nats.Msg) error) error
	Publish(topic string, payload []byte) error
	Request(topic string, payload []byte) (*nats.Msg, error)
}

type conn struct {
	servers []string
	token   string
	group   string
	name    string

	*nats.Conn
	mx   *sync.RWMutex
	subs map[string]*nats.Subscription
}

func NewConn(cfg *ConnConfig) (driver.Conn, error) {
	conn := &conn{
		servers: cfg.Servers,
		token:   cfg.Token,
		group:   cfg.Group,
		name:    cfg.Name,

		mx:   &sync.RWMutex{},
		subs: make(map[string]*nats.Subscription),
	}

	var err error

	conn.Conn, err = nats.Connect(strings.Join(conn.servers, ", "), nats.Token(conn.token), nats.Name(conn.name))

	if err != nil {
		return nil, errtpl.ErrConnect(err, "nats")
	}

	return conn, nil
}

func (c *conn) Serve(ctx context.Context) error {
	<-ctx.Done()

	for topic := range c.subs {
		if err := c.unsubscribe(topic); err != nil {
			log.Error(err)
		}
	}

	return nil
}

func (c *conn) unsubscribe(topic string) error {
	c.mx.Lock()
	defer c.mx.Unlock()

	sub, ok := c.subs[topic]

	if !ok {
		return nil
	}

	if err := sub.Unsubscribe(); err != nil {
		return msgbroker.ErrUnsubscribe(prepareError(err), topic)
	}

	delete(c.subs, topic)

	msgbroker.LogDebugUnsubscribed(topic)

	return nil
}

func (c *conn) Close() error {
	c.Conn.Close()

	return nil
}

func (c *conn) Receiver(route *entity.Route) driver.Receiver {
	return &receiver{
		conn:  c,
		topic: route.Topic,
	}
}

func (c *conn) Sender(route *entity.Route) driver.Sender {
	return &sender{
		conn:  c,
		topic: route.Topic,
	}
}

func (c *conn) Subscribe(topic string, handler func(*nats.Msg) error) error {
	// Check if already subscribed
	c.mx.RLock()
	_, ok := c.subs[topic] //nolint:ifshort
	c.mx.RUnlock()

	if ok {
		return nil
	}

	sub, err := c.QueueSubscribe(topic, c.group, func(msg *nats.Msg) {
		if err := handler(msg); err != nil {
			msgbroker.LogErrorHandle(err, topic)
		}
	})

	if err != nil {
		return msgbroker.ErrSubscribe(prepareError(err), topic)
	}

	c.mx.Lock()
	defer c.mx.Unlock()

	c.subs[topic] = sub

	msgbroker.LogDebugSubscribed(topic)

	return nil
}

func (c *conn) Publish(topic string, payload []byte) error {
	err := c.Conn.Publish(topic, payload)

	if err != nil {
		return msgbroker.ErrPublish(prepareError(err), topic)
	}

	msgbroker.LogDebugPublished(topic, payload)

	return nil
}

func (c *conn) Request(topic string, payload []byte) (*nats.Msg, error) {
	msg, err := c.Conn.Request(topic, payload, requestTimeout)

	if err != nil {
		return nil, msgbroker.ErrBadReply(prepareError(err), topic)
	}

	msgbroker.LogDebugRequested(topic, nil, nil)

	return msg, nil
}

func prepareError(err error) error {
	if errors.Is(err, nats.ErrNoResponders) {
		return msgbroker.ErrNoResponders
	}

	return err
}
