package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"NATter/batcher"
	"NATter/driver"
	"NATter/entity"
	m "NATter/mock"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewRouter(t *testing.T) {
	connBroker := &m.DriverConn{}
	connHTTP := &m.DriverConn{}

	receiverHTTP := &m.DriverReceiver{}
	receiverBroker := &m.DriverReceiver{}

	senderBroker := &m.DriverSender{}

	connBroker.
		On("Sender", &entity.Route{
			Mode: entity.RouteMode("http-broker-oneway"),
		}).
		Return(senderBroker)
	connBroker.
		On("Receiver", &entity.Route{
			Mode: entity.RouteMode("broker-http-twoway"),
			Batching: &entity.RouteBatching{
				Timeout:  30,
				Capacity: 5,
			},
		}).
		Return(receiverBroker)

	connHTTP.
		On("Receiver", &entity.Route{
			Mode: entity.RouteMode("http-broker-oneway"),
		}).
		Return(receiverHTTP)
	connHTTP.
		On("Sender", &entity.Route{
			Mode: entity.RouteMode("broker-http-twoway"),
			Batching: &entity.RouteBatching{
				Timeout:  30,
				Capacity: 5,
			},
		}).
		Return(&m.DriverSender{})

	receiverHTTP.On("Listen", senderBroker).Return(nil)

	receiverBroker.On("ListenRequest", mock.AnythingOfType("*batcher.batcher")).Return(nil)

	router, err := NewRouter(&RouterConfig{
		Routes: []*entity.Route{
			{
				Mode: entity.RouteMode("http-broker-oneway"),
			},
			{
				Mode: entity.RouteMode("broker-http-twoway"),
				Batching: &entity.RouteBatching{
					Timeout:  30,
					Capacity: 5,
				},
			},
		},
	}, map[string]driver.Conn{
		"broker": connBroker,
		"http":   connHTTP,
	})

	assert.Nil(t, err)
	assert.NotNil(t, router)

	connBroker.AssertExpectations(t)
	connHTTP.AssertExpectations(t)
	receiverHTTP.AssertExpectations(t)
	receiverBroker.AssertExpectations(t)
}

func TestNewRouterOnUnknownReceiverDriverConn(t *testing.T) {
	router, err := NewRouter(&RouterConfig{
		Routes: []*entity.Route{
			{
				Mode: entity.RouteMode("UNKNOWN-broker-oneway"),
			},
		},
	}, map[string]driver.Conn{
		"broker": &m.DriverConn{},
		"http":   &m.DriverConn{},
	})

	assert.Error(t, err)
	assert.Nil(t, router)
}

func TestNewRouterOnUnknownSenderDriverConn(t *testing.T) {
	router, err := NewRouter(&RouterConfig{
		Routes: []*entity.Route{
			{
				Mode: entity.RouteMode("http-UNKNOWN-oneway"),
			},
		},
	}, map[string]driver.Conn{
		"broker": &m.DriverConn{},
		"http":   &m.DriverConn{},
	})

	assert.Error(t, err)
	assert.Nil(t, router)
}

func TestNewRouterOnNewBatcherError(t *testing.T) {
	connHTTP := &m.DriverConn{}

	connHTTP.
		On("Sender", &entity.Route{
			Mode: entity.RouteMode("broker-http-twoway"),
			Batching: &entity.RouteBatching{
				Timeout:  0,
				Capacity: 0,
			},
		}).
		Return(&m.DriverSender{})

	router, err := NewRouter(&RouterConfig{
		Routes: []*entity.Route{
			{
				Mode: entity.RouteMode("broker-http-twoway"),
				Batching: &entity.RouteBatching{
					Timeout:  0,
					Capacity: 0,
				},
			},
		},
	}, map[string]driver.Conn{
		"broker": &m.DriverConn{},
		"http":   connHTTP,
	})

	assert.Error(t, err)
	assert.Nil(t, router)
}

func TestNewRouterOnUnknownDirection(t *testing.T) {
	DriverConnBroker := &m.DriverConn{}

	DriverConnBroker.
		On("Sender", &entity.Route{
			Mode: entity.RouteMode("http-broker-UNKNOWN"),
		}).
		Return(&m.DriverSender{})

	router, err := NewRouter(&RouterConfig{
		Routes: []*entity.Route{
			{
				Mode: entity.RouteMode("http-broker-UNKNOWN"),
			},
		},
	}, map[string]driver.Conn{
		"broker": DriverConnBroker,
		"http":   &m.DriverConn{},
	})

	assert.Error(t, err)
	assert.Nil(t, router)

	DriverConnBroker.AssertExpectations(t)
}

func TestRouterRun(t *testing.T) {
	conn1 := &m.DriverConn{}
	conn2 := &m.DriverConn{}
	bat1 := &m.Batcher{}
	bat2 := &m.Batcher{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	conn1.On("Serve", ctx).Return(nil)
	conn2.On("Serve", ctx).Return(nil)

	bat1.On("Run", ctx)
	bat2.On("Run", ctx)

	router := &Router{
		conns: map[string]driver.Conn{
			"conn1": conn1,
			"conn2": conn2,
		},
		batchers: []batcher.Batcher{
			bat1,
			bat2,
		},
		wg: &sync.WaitGroup{},
	}

	router.Run(ctx)

	conn1.AssertExpectations(t)
	conn2.AssertExpectations(t)
	bat1.AssertExpectations(t)
	bat2.AssertExpectations(t)
}

func TestRouterRunOnServeError(t *testing.T) {
	conn := &m.DriverConn{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	conn.On("Serve", ctx).Return(errors.New("error"))

	router := &Router{
		conns: map[string]driver.Conn{
			"conn": conn,
		},
		wg: &sync.WaitGroup{},
	}

	router.Run(ctx)

	conn.AssertExpectations(t)
}
