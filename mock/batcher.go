package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type Batcher struct {
	mock.Mock
}

func (b *Batcher) Run(ctx context.Context) {
	b.Called(ctx)
}

func (b *Batcher) Send(msg []byte) error {
	args := b.Called(msg)

	return args.Error(0)
}

func (b *Batcher) Request(msg []byte) ([]byte, error) {
	args := b.Called(msg)

	return args.Get(0).([]byte), args.Error(1)
}

type BatcherEncoder struct {
	mock.Mock
}

func (e *BatcherEncoder) Marshal(msgs [][]byte) ([]byte, error) {
	args := e.Called(msgs)

	return args.Get(0).([]byte), args.Error(1)
}
