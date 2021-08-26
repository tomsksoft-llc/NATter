package batcher

import (
	"context"
	"sync"
	"time"

	"NATter/batcher/encoder"
	"NATter/driver"
	"NATter/log"

	"github.com/pkg/errors"
)

type Config struct {
	Timeout  uint32 // seconds
	Capacity uint32
}

type Batcher interface {
	Run(context.Context)
	Send(msg []byte) error
	Request(msg []byte) ([]byte, error)
}

type batcher struct {
	sender driver.Sender
	enc    encoder.Encoder

	msgChan chan []byte
	wg      *sync.WaitGroup

	timeout  time.Duration
	capacity uint32
}

func New(cfg *Config, sender driver.Sender, enc encoder.Encoder) (Batcher, error) {
	batcher := &batcher{
		sender:   sender,
		enc:      enc,
		msgChan:  make(chan []byte),
		wg:       &sync.WaitGroup{},
		timeout:  time.Duration(cfg.Timeout) * time.Second,
		capacity: cfg.Capacity,
	}

	if cfg.Timeout == 0 && cfg.Capacity == 0 {
		return nil, errors.New("both timeout and capacity have default 0 value")
	}

	return batcher, nil
}

func (b *batcher) Run(ctx context.Context) {
	ticker := b.prepareTicker()
	msgs := [][]byte{}

OUTER:
	for {
		select {
		case msg := <-b.msgChan:
			msgs = append(msgs, msg)

			log.Debug("pushed new message to batch")

			if len(msgs) < int(b.capacity) || b.capacity == 0 {
				break
			}

			b.releaseBatch(msgs)

			msgs = nil

			if b.timeout > 0 {
				ticker.Reset(b.timeout)
			}
		case <-ticker.C:
			b.releaseBatch(msgs)

			msgs = nil
		case <-ctx.Done():
			b.releaseBatch(msgs)

			break OUTER
		}
	}

	b.wg.Wait()
}

func (b *batcher) prepareTicker() (t *time.Ticker) {
	if b.timeout > 0 {
		t = time.NewTicker(b.timeout)
	} else {
		t = time.NewTicker(1)
		t.Stop()
	}

	return t
}

func (b *batcher) releaseBatch(msgs [][]byte) {
	if len(msgs) == 0 {
		return
	}

	b.wg.Add(1)

	go func() {
		defer b.wg.Done()

		batch, err := b.enc.Marshal(msgs)

		if err != nil {
			log.Error(err)

			return
		}

		if err := b.sender.Send(batch); err != nil {
			log.Error(err)

			return
		}

		log.Debugf("released new batch of %d messages", len(msgs))
	}()
}

func (b *batcher) Send(msg []byte) error {
	b.msgChan <- msg

	return nil
}

func (b *batcher) Request(msg []byte) ([]byte, error) {
	return nil, errors.New("request is not supported in protobuf batching")
}
