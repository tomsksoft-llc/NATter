package nats

import (
	"context"
	"sync"
	"testing"
	"time"

	"NATter/driver/msgbroker"
	"NATter/entity"

	"github.com/nats-io/nats-server/v2/server"
	servertest "github.com/nats-io/nats-server/v2/test"
	nats "github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConnTestSuite struct {
	suite.Suite
	srv  *server.Server
	conn *conn
}

func (s *ConnTestSuite) SetupSuite() {
	opts := servertest.DefaultTestOptions
	opts.Port = 3000
	s.srv = servertest.RunServer(&opts)
}

func (s *ConnTestSuite) SetupTest() {
	s.conn = &conn{
		mx:   &sync.RWMutex{},
		subs: make(map[string]*nats.Subscription),
	}

	var err error

	s.conn.Conn, err = nats.Connect("nats://localhost:3000", nats.Token(""), nats.Name(""))

	assert.Nil(s.T(), err)
}

func (s *ConnTestSuite) TearDownTest() {
	s.conn.Close()
}

func (s *ConnTestSuite) TearDownSuite() {
	s.srv.Shutdown()
	s.srv.WaitForShutdown()
}

func (s *ConnTestSuite) TestNewConn() {
	conn, err := NewConn(&ConnConfig{
		Servers: []string{"nats://localhost:3000"},
	})

	assert.NotNil(s.T(), conn)
	assert.Nil(s.T(), err)
}

func (s *ConnTestSuite) TestNewConnOnConnectError() {
	conn, err := NewConn(&ConnConfig{
		Servers: []string{"nats://localhost:1"},
		Token:   "random",
	})

	assert.Nil(s.T(), conn)
	assert.Error(s.T(), err)
}

func (s *ConnTestSuite) TestServe() {
	err := s.conn.Subscribe("hey", func(*nats.Msg) error {
		return nil
	})

	assert.Nil(s.T(), err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	err = s.conn.Serve(ctx)

	<-ctx.Done()

	assert.Nil(s.T(), err)
	assert.Zero(s.T(), s.conn.NumSubscriptions())
}

func (s *ConnTestSuite) TestServeOnUnsubscribeError() {
	err := s.conn.Subscribe("hey", func(*nats.Msg) error {
		return nil
	})

	assert.Nil(s.T(), err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	s.TearDownTest()

	err = s.conn.Serve(ctx)

	<-ctx.Done()

	assert.Nil(s.T(), err)
}

func (s *ConnTestSuite) TestUnsubscribe() {
	err := s.conn.Subscribe("hey", func(*nats.Msg) error {
		return nil
	})

	assert.Nil(s.T(), err)

	err = s.conn.unsubscribe("hey")

	assert.Nil(s.T(), err)
	assert.Zero(s.T(), s.conn.NumSubscriptions())
}

func (s *ConnTestSuite) TestUnsubscribeOnError() {
	err := s.conn.Subscribe("hey", func(*nats.Msg) error {
		return nil
	})

	assert.Nil(s.T(), err)

	s.TearDownTest()

	err = s.conn.unsubscribe("hey")

	assert.Error(s.T(), err)
}

func (s *ConnTestSuite) TestUnsubscribeTwice() {
	err := s.conn.Subscribe("hey", func(*nats.Msg) error {
		return nil
	})

	assert.Nil(s.T(), err)

	err = s.conn.unsubscribe("hey")

	assert.Nil(s.T(), err)
	assert.Zero(s.T(), s.conn.NumSubscriptions())

	err = s.conn.unsubscribe("hey")

	assert.Nil(s.T(), err)
	assert.Zero(s.T(), s.conn.NumSubscriptions())
}

func (s *ConnTestSuite) TestReceiver() {
	rec, ok := s.conn.Receiver(&entity.Route{
		Topic: "topic",
	}).(*receiver)

	if !ok {
		s.Fail("type assertion error")
	}

	assert.Equal(s.T(), "topic", rec.topic)
	assert.Equal(s.T(), s.conn, rec.conn)
}

func (s *ConnTestSuite) TestSender() {
	snd, ok := s.conn.Sender(&entity.Route{
		Topic: "topic",
	}).(*sender)

	if !ok {
		s.Fail("type assertion error")
	}

	assert.Equal(s.T(), "topic", snd.topic)
	assert.Equal(s.T(), s.conn, snd.conn)
}

func (s *ConnTestSuite) TestSubscribe() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	err := s.conn.Subscribe("hey", func(msg *nats.Msg) error {
		assert.Equal(s.T(), "hey", msg.Subject)
		assert.Equal(s.T(), []byte("hello"), msg.Data)

		cancel()

		return nil
	})

	assert.Nil(s.T(), err)

	err = s.conn.Conn.Publish("hey", []byte("hello"))

	<-ctx.Done()

	assert.Equal(s.T(), context.Canceled, ctx.Err())
	assert.Nil(s.T(), err)
}

func (s *ConnTestSuite) TestSubscribeOnError() {
	s.TearDownTest()

	err := s.conn.Subscribe("hey", func(*nats.Msg) error {
		return nil
	})

	assert.Error(s.T(), err)
}

func (s *ConnTestSuite) TestSubscribeOnSecondTime() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	err := s.conn.Subscribe("hey", func(msg *nats.Msg) error {
		assert.Equal(s.T(), "hey", msg.Subject)
		assert.Equal(s.T(), []byte("hello"), msg.Data)

		cancel()

		return nil
	})

	assert.Nil(s.T(), err)

	err = s.conn.Subscribe("hey", func(msg *nats.Msg) error {
		assert.Equal(s.T(), "hey", msg.Subject)
		assert.Equal(s.T(), []byte("h0lle"), msg.Data)

		cancel()

		return nil
	})

	assert.Nil(s.T(), err)

	err = s.conn.Conn.Publish("hey", []byte("hello"))

	<-ctx.Done()

	assert.Equal(s.T(), context.Canceled, ctx.Err())
	assert.Nil(s.T(), err)
}

func (s *ConnTestSuite) TestSubscribeOnHandlerError() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	err := s.conn.Subscribe("hey", func(msg *nats.Msg) error {
		assert.Equal(s.T(), "hey", msg.Subject)
		assert.Equal(s.T(), []byte("hello"), msg.Data)

		cancel()

		return errors.New("error")
	})

	assert.Nil(s.T(), err)

	err = s.conn.Conn.Publish("hey", []byte("hello"))

	<-ctx.Done()

	assert.Equal(s.T(), context.Canceled, ctx.Err())
	assert.Nil(s.T(), err)
}

func (s *ConnTestSuite) TestPublish() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	conn, err := s.conn.Conn.Subscribe("hey", func(msg *nats.Msg) {
		assert.EqualValues(s.T(), "hello", msg.Data)

		cancel()
	})

	assert.Nil(s.T(), err)

	err = s.conn.Publish("hey", []byte("hello"))

	<-ctx.Done()

	assert.Equal(s.T(), context.Canceled, ctx.Err())
	assert.Nil(s.T(), err)

	err = conn.Unsubscribe()
	assert.Nil(s.T(), err)
}

func (s *ConnTestSuite) TestPublishOnError() {
	s.TearDownTest()

	err := s.conn.Publish("hey", []byte("hello"))

	assert.Error(s.T(), err)
}

func (s *ConnTestSuite) TestRequest() {
	_, err := s.conn.Conn.Subscribe("hey", func(m *nats.Msg) {
		err := m.Respond([]byte("hi"))

		assert.Nil(s.T(), err)
	})

	assert.Nil(s.T(), err)

	msg, err := s.conn.Request("hey", []byte("hello"))

	assert.Equal(s.T(), []byte("hi"), msg.Data)
	assert.Nil(s.T(), err)
}

func (s *ConnTestSuite) TestRequestOnNoSubscribers() {
	msg, err := s.conn.Request("hey", []byte("hello"))

	assert.Nil(s.T(), msg)
	assert.True(s.T(), errors.Is(err, msgbroker.ErrNoResponders))
}

func (s *ConnTestSuite) TestRequestOnError() {
	s.TearDownTest()

	msg, err := s.conn.Request("hey", []byte("hello"))

	assert.Nil(s.T(), msg)
	assert.Error(s.T(), err)
}

func TestMsgBroker(t *testing.T) {
	suite.Run(t, &ConnTestSuite{})
}
