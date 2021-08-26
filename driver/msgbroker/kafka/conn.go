package kafka

import (
	"context"
	"strings"

	"NATter/driver"
	"NATter/driver/msgbroker"
	"NATter/entity"
	"NATter/errtpl"
	"NATter/log"

	"github.com/Shopify/sarama"
)

const DriverName = "kafka"

type ConnConfig struct {
	Version string
	Servers []string
	Group   string
}

type Conn interface {
	Subscribe(topic string, handler func([]byte) error) error
	Publish(topic string, payload []byte) error
}

type conn struct {
	version string
	servers []string
	group   string

	consumer sarama.ConsumerGroup
	producer sarama.AsyncProducer

	handlers map[string]func([]byte) error
	gch      *consumerHandler
}

func NewConn(cfg *ConnConfig) (driver.Conn, error) {
	conn := &conn{
		version: cfg.Version,
		servers: cfg.Servers,
		group:   cfg.Group,

		handlers: map[string]func([]byte) error{},
		gch:      &consumerHandler{},
	}

	saramaConf := sarama.NewConfig()
	saramaConf.Consumer.Fetch.Default = 1024 * 500

	var err error

	saramaConf.Version, err = sarama.ParseKafkaVersion(conn.version)

	if err != nil {
		return nil, err
	}

	conn.producer, err = sarama.NewAsyncProducer(conn.servers, saramaConf)

	if err != nil {
		return nil, errtpl.ErrConnect(err, "kafka")
	}

	conn.consumer, err = sarama.NewConsumerGroup(conn.servers, cfg.Group, saramaConf)

	if err != nil {
		return nil, errtpl.ErrConnect(err, "kafka")
	}

	return conn, nil
}

func (c *conn) Serve(ctx context.Context) error {
	topics := []string{}

	for topic := range c.handlers {
		topics = append(topics, topic)
	}

	c.gch.reset(c.handlers)

	go func() {
		for {
			if err := c.consumer.Consume(ctx, topics, c.gch); err != nil {
				log.Error(msgbroker.ErrSubscribe(err, strings.Join(topics, ", ")))
			}

			if ctx.Err() != nil {
				return
			}

			c.gch.ready = make(chan bool)
		}
	}()

	<-c.gch.ready

	return nil
}

func (c *conn) Close() error {
	c.consumer.Close()
	c.producer.Close()

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

func (c *conn) Publish(topic string, message []byte) error {
	c.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	msgbroker.LogDebugPublished(topic, nil)

	return nil
}

func (c *conn) Subscribe(topic string, handler func([]byte) error) error {
	c.handlers[topic] = handler

	msgbroker.LogDebugSubscribed(topic)

	return nil
}

type consumerHandler struct {
	handlers map[string]func([]byte) error
	ready    chan bool
}

func (ch *consumerHandler) reset(handlers map[string]func([]byte) error) {
	ch.handlers = handlers
	ch.ready = make(chan bool)
}

func (ch consumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(ch.ready)

	return nil
}

func (ch consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (ch consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		msgbroker.LogDebugReceived(msg.Topic, nil)

		err := ch.handlers[msg.Topic](msg.Value)

		if err != nil {
			msgbroker.LogErrorHandle(err, msg.Topic)

			continue
		}

		sess.MarkMessage(msg, "")
	}

	return nil
}
