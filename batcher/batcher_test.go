package batcher

import (
	"context"
	"testing"
	"time"

	m "NATter/mock"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewOnDefaultParameters(t *testing.T) {
	bat, err := New(&Config{}, &m.DriverSender{}, &m.BatcherEncoder{})

	assert.Error(t, err)
	assert.Nil(t, bat)
}

func TestBatcherRunOnFullBatch(t *testing.T) {
	sender := &m.DriverSender{}
	enc := &m.BatcherEncoder{}

	sender.
		On("Send", []byte("batch-of-data")).
		Return(nil).Once()

	enc.
		On("Marshal", [][]byte{
			[]byte("some-data"),
			[]byte("some-data"),
			[]byte("some-data"),
		}).
		Return([]byte("batch-of-data"), nil).Once()

	bat, err := New(&Config{
		Timeout:  30,
		Capacity: 3,
	}, sender, enc)

	assert.Nil(t, err)
	assert.NotNil(t, bat)

	go func() {
		for i := 0; i < 3; i++ {
			err := bat.Send([]byte("some-data"))

			assert.Nil(t, err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	bat.Run(ctx)

	sender.AssertExpectations(t)
	enc.AssertExpectations(t)
}

func TestBatcherRunOnTimeout(t *testing.T) {
	sender := &m.DriverSender{}
	enc := &m.BatcherEncoder{}

	sender.
		On("Send", []byte("batch-of-data")).
		Return(nil).Once()

	enc.
		On("Marshal", [][]byte{
			[]byte("some-data"),
		}).
		Return([]byte("batch-of-data"), nil).Once()

	bat, err := New(&Config{
		Timeout:  1,
		Capacity: 3,
	}, sender, enc)

	assert.Nil(t, err)
	assert.NotNil(t, bat)

	go func() {
		err := bat.Send([]byte("some-data"))

		assert.Nil(t, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1001)
	defer cancel()

	bat.Run(ctx)

	sender.AssertExpectations(t)
	enc.AssertExpectations(t)
}

func TestBatcherRunOnMarshalError(t *testing.T) {
	enc := &m.BatcherEncoder{}

	enc.
		On("Marshal", [][]byte{
			[]byte("some-data"),
		}).
		Return([]byte(nil), errors.New("error")).Once()

	bat, err := New(&Config{
		Timeout:  30,
		Capacity: 3,
	}, &m.DriverSender{}, enc)

	assert.Nil(t, err)
	assert.NotNil(t, bat)

	go func() {
		err := bat.Send([]byte("some-data"))

		assert.Nil(t, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	bat.Run(ctx)

	enc.AssertExpectations(t)
}

func TestBatcherRunOnSendError(t *testing.T) {
	sender := &m.DriverSender{}
	enc := &m.BatcherEncoder{}

	sender.
		On("Send", []byte("batch-of-data")).
		Return(errors.New("error"))

	enc.
		On("Marshal", [][]byte{
			[]byte("some-data"),
		}).
		Return([]byte("batch-of-data"), nil).Once()

	bat, err := New(&Config{
		Timeout:  30,
		Capacity: 1,
	}, sender, enc)

	assert.Nil(t, err)
	assert.NotNil(t, bat)

	go func() {
		err := bat.Send([]byte("some-data"))

		assert.Nil(t, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	bat.Run(ctx)

	sender.AssertExpectations(t)
	enc.AssertExpectations(t)
}

func TestBatcherRunOnZeroTimeout(t *testing.T) {
	sender := &m.DriverSender{}

	bat, err := New(&Config{
		Timeout:  0,
		Capacity: 3,
	}, sender, &m.BatcherEncoder{})

	assert.Nil(t, err)
	assert.NotNil(t, bat)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	bat.Run(ctx)
}

func TestBatcherSend(t *testing.T) {
	sender := &m.DriverSender{}
	enc := &m.BatcherEncoder{}

	sender.
		On("Send", []byte("batch-of-data")).
		Return(nil).Once()

	enc.
		On("Marshal", [][]byte{
			[]byte("some-data"),
		}).
		Return([]byte("batch-of-data"), nil).Once()

	bat, err := New(&Config{
		Timeout:  30,
		Capacity: 3,
	}, sender, enc)

	assert.Nil(t, err)
	assert.NotNil(t, bat)

	go func() {
		err := bat.Send([]byte("some-data"))

		assert.Nil(t, err)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	bat.Run(ctx)

	sender.AssertExpectations(t)
	enc.AssertExpectations(t)
}

func TestBatcherRequest(t *testing.T) {
	bat, err := New(&Config{
		Timeout:  30,
		Capacity: 3,
	}, &m.DriverSender{}, &m.BatcherEncoder{})

	assert.Nil(t, err)
	assert.NotNil(t, bat)

	respb, err := bat.Request([]byte("some-data"))

	assert.Error(t, err)
	assert.Nil(t, respb)
}
