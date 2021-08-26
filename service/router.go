package service

import (
	"context"
	"sync"

	"NATter/batcher"
	"NATter/batcher/encoder/protobuf"
	"NATter/driver"
	"NATter/entity"
	"NATter/log"

	"github.com/pkg/errors"
)

var (
	ErrUnknownConn = errors.New("unknown connection")
)

type RouterConfig struct {
	Routes []*entity.Route
}

type Router struct {
	conns    map[string]driver.Conn
	batchers []batcher.Batcher

	wg *sync.WaitGroup
}

func NewRouter(cfg *RouterConfig, conns map[string]driver.Conn) (*Router, error) {
	router := &Router{
		conns:    conns,
		batchers: []batcher.Batcher{},
		wg:       &sync.WaitGroup{},
	}

	if err := router.registerRoutes(cfg.Routes); err != nil {
		return nil, err
	}

	return router, nil
}

func (router *Router) registerRoutes(routes []*entity.Route) (err error) {
	for _, r := range routes {
		modeComp := r.Mode.Components()

		receiverConn, ok := router.conns[modeComp.Receiver]

		if !ok {
			return errors.Wrap(ErrUnknownConn, modeComp.Receiver)
		}

		senderConn, ok := router.conns[modeComp.Sender]

		if !ok {
			return errors.Wrap(ErrUnknownConn, modeComp.Sender)
		}

		sender := senderConn.Sender(r)

		if r.Batching != nil {
			bat, err := batcher.New(&batcher.Config{
				Timeout:  r.Batching.Timeout,
				Capacity: r.Batching.Capacity,
			}, sender, &protobuf.Encoder{})

			if err != nil {
				return err
			}

			router.batchers = append(router.batchers, bat)

			sender = bat
		}

		switch modeComp.Direction {
		case entity.RouteDirectionOneway:
			err = receiverConn.Receiver(r).Listen(sender)
		case entity.RouteDirectionTwoway:
			err = receiverConn.Receiver(r).ListenRequest(sender)
		default:
			err = errors.Errorf("unknown route direction: %s", modeComp.Direction)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (router *Router) Run(ctx context.Context) {
	for _, bat := range router.batchers {
		router.wg.Add(1)

		go func(bat batcher.Batcher) {
			defer router.wg.Done()

			bat.Run(ctx)
		}(bat)
	}

	for _, conn := range router.conns {
		router.wg.Add(1)

		go func(conn driver.Conn) {
			defer router.wg.Done()

			if err := conn.Serve(ctx); err != nil {
				log.Error(err)
			}
		}(conn)
	}

	router.wg.Wait()
}
