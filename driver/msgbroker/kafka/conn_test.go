package kafka

import (
	"context"
	"testing"
	"time"

	"NATter/entity"
	m "NATter/mock"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func testConnEnv(t *testing.T) (*m.ConsumerGroup, *mocks.AsyncProducer, *conn) {
	t.Helper()

	producer := mocks.NewAsyncProducer(t, sarama.NewConfig())
	consumer := &m.ConsumerGroup{}

	conn := &conn{
		consumer: consumer,
		producer: producer,
		handlers: map[string]func([]byte) error{},
		gch:      &consumerHandler{},
	}

	return consumer, producer, conn
}

func TestNewConnOnError(t *testing.T) {
	_, err := NewConn(&ConnConfig{})

	assert.NotNil(t, err)
}

func TestConnServe(t *testing.T) {
	consumer, _, conn := testConnEnv(t)

	err := conn.Subscribe("topic", func([]byte) error { return nil })
	assert.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err = conn.Serve(ctx)

		assert.Nil(t, err)
	}()

	gch := conn.gch

	consumer.
		On("Consume", ctx, []string{"topic"}, gch).
		Return(nil)

	time.Sleep(time.Millisecond * 4)

	consumer.AssertExpectations(t)
}

func TestConnServeOnError(t *testing.T) {
	consumer, _, conn := testConnEnv(t)

	err := conn.Subscribe("topic", func([]byte) error { return nil })
	assert.Nil(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err = conn.Serve(ctx)

		assert.Nil(t, err)
	}()

	gch := conn.gch

	consumer.
		On("Consume", ctx, []string{"topic"}, gch).
		Return(errors.New("error"))

	time.Sleep(time.Microsecond * 100)

	consumer.AssertExpectations(t)
}

func TestConnClose(t *testing.T) {
	consumer, _, conn := testConnEnv(t)

	consumer.
		On("Close").Once().
		Return(nil)

	conn.Close()

	consumer.AssertExpectations(t)
}

func TestConnReceiver(t *testing.T) {
	_, _, conn := testConnEnv(t)

	rec, ok := conn.Receiver(&entity.Route{
		Topic: "topic",
	}).(*receiver)

	if !ok {
		assert.Fail(t, "type assertion error")
	}

	assert.Equal(t, "topic", rec.topic)
	assert.Equal(t, conn, rec.conn)
}

func TestConnSender(t *testing.T) {
	_, _, conn := testConnEnv(t)

	snd, ok := conn.Sender(&entity.Route{
		Topic: "topic",
	}).(*sender)

	if !ok {
		assert.Fail(t, "type assertion error")
	}

	assert.Equal(t, "topic", snd.topic)
	assert.Equal(t, conn, snd.conn)
}

func TestConnPublish(t *testing.T) {
	_, producer, conn := testConnEnv(t)

	producer.ExpectInputAndSucceed()

	err := conn.Publish("topic", []byte("message"))

	assert.Nil(t, err)
}

func TestConnSubscribe(t *testing.T) {
	_, _, conn := testConnEnv(t)

	err := conn.Subscribe("topic", func([]byte) error { return nil })
	assert.Nil(t, err)
}

func TestConsumerHandlerSetup(t *testing.T) {
	gch := consumerHandler{}

	gch.reset(map[string]func([]byte) error{})

	err := gch.Setup(&m.ConsumerGroupSession{})

	assert.Nil(t, err)
}

func TestConsumerHandlerCleanup(t *testing.T) {
	gch := consumerHandler{}

	err := gch.Cleanup(&m.ConsumerGroupSession{})

	assert.Nil(t, err)
}

func TestConsumerHandlerConsumeClaim(t *testing.T) {
	sess := &m.ConsumerGroupSession{}
	claim := &m.ConsumerGroupClaim{}

	msg := &sarama.ConsumerMessage{Topic: "internal"}

	claim.
		On("Messages").Once().
		Return(msg)

	sess.
		On("MarkMessage", msg, "").Once().
		Return()

	gch := consumerHandler{}

	gch.reset(map[string]func([]byte) error{"internal": func([]byte) error { return nil }})

	err := gch.ConsumeClaim(sess, claim)

	assert.Nil(t, err)
	sess.AssertExpectations(t)
	claim.AssertExpectations(t)
}

func TestConsumerHandlerConsumeClaimOnError(t *testing.T) {
	sess := &m.ConsumerGroupSession{}
	claim := &m.ConsumerGroupClaim{}

	msg := &sarama.ConsumerMessage{Topic: "internal"}

	claim.
		On("Messages").Once().
		Return(msg)

	gch := consumerHandler{}

	gch.reset(map[string]func([]byte) error{
		"internal": func([]byte) error { return errors.New("error") },
	})

	err := gch.ConsumeClaim(sess, claim)

	assert.Nil(t, err)
	sess.AssertExpectations(t)
	claim.AssertExpectations(t)
}
