package nats

import (
	"sync"
	"testing"
	"time"

	m "NATter/mock"

	"github.com/nats-io/nats-server/v2/server"
	servertest "github.com/nats-io/nats-server/v2/test"
	nats "github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestReceiverListen(t *testing.T) {
	conn := &m.DriverNatsConn{}

	conn.
		On("Subscribe", "topic", mock.AnythingOfType("func(*nats.Msg) error")).
		Run(func(args mock.Arguments) {
			msg := &nats.Msg{
				Data: []byte("some-data"),
			}

			err := args.Get(1).(func(*nats.Msg) error)(msg)

			assert.Nil(t, err)
		}).
		Return(nil)

	sender := &m.DriverSender{}

	sender.
		On("Send", []byte("some-data")).
		Return(nil)

	receiver := &receiver{
		conn:  conn,
		topic: "topic",
	}

	err := receiver.Listen(sender)

	assert.Nil(t, err)
}

type ReceiverTestSuite struct {
	suite.Suite
	srv  *server.Server
	conn *conn
}

func (s *ReceiverTestSuite) SetupSuite() {
	opts := servertest.DefaultTestOptions
	opts.Port = 3000
	s.srv = servertest.RunServer(&opts)
}

func (s *ReceiverTestSuite) SetupTest() {
	s.conn = &conn{
		mx:   &sync.RWMutex{},
		subs: make(map[string]*nats.Subscription),
	}

	var err error

	s.conn.Conn, err = nats.Connect("nats://localhost:3000", nats.Token(""), nats.Name(""))

	assert.Nil(s.T(), err)
}

func (s *ReceiverTestSuite) TearDownTest() {
	s.conn.Close()
}

func (s *ReceiverTestSuite) TearDownSuite() {
	s.srv.Shutdown()
	s.srv.WaitForShutdown()
}

func (s *ReceiverTestSuite) TestListenRequest() {
	sender := &m.DriverSender{}

	sender.
		On("Request", []byte("request-data")).
		Return([]byte("response-data"), nil)

	receiver := &receiver{
		conn:  s.conn,
		topic: "topic",
	}

	err := receiver.ListenRequest(sender)

	assert.Nil(s.T(), err)

	msg, err := s.conn.Request("topic", []byte("request-data"))

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), []byte("response-data"), msg.Data)
}

func (s *ReceiverTestSuite) TestListenRequestOnError() {
	sender := &m.DriverSender{}

	sender.
		On("Request", []byte("request-data")).
		Return([]byte(nil), errors.New("error"))

	receiver := &receiver{
		conn:  s.conn,
		topic: "topic",
	}

	err := receiver.ListenRequest(sender)

	assert.Nil(s.T(), err)

	msg, err := s.conn.Request("topic", []byte("request-data"))

	assert.Error(s.T(), err)
	assert.Nil(s.T(), msg)
}

func (s *ReceiverTestSuite) TestListenRequestOnRespondError() {
	sender := &m.DriverSender{}

	sender.
		On("Request", []byte("request-data")).
		Return([]byte("response-data"), nil).
		After(time.Millisecond * 5)

	receiver := &receiver{
		conn:  s.conn,
		topic: "topic",
	}

	err := receiver.ListenRequest(sender)

	assert.Nil(s.T(), err)

	go func() {
		msg, err := s.conn.Request("topic", []byte("request-data"))

		assert.Error(s.T(), err)
		assert.Nil(s.T(), msg)
	}()

	time.Sleep(time.Millisecond)

	s.TearDownTest()

	time.Sleep(time.Millisecond * 10)
}

func TestReceiver(t *testing.T) {
	suite.Run(t, &ReceiverTestSuite{})
}
