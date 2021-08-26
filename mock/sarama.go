package mock

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/mock"
)

type ConsumerGroup struct {
	mock.Mock
}

func (cg *ConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	args := cg.Called(ctx, topics, handler)

	return args.Error(0)
}

func (cg *ConsumerGroup) Errors() <-chan error {
	args := cg.Called()

	out := make(chan error)

	go func() {
		out <- args.Error(0)
	}()

	return out
}

func (cg *ConsumerGroup) Close() error {
	args := cg.Called()

	return args.Error(0)
}

type ConsumerGroupSession struct {
	mock.Mock
}

func (cgs *ConsumerGroupSession) Claims() map[string][]int32 {
	args := cgs.Called()

	return args.Get(0).(map[string][]int32)
}

func (cgs *ConsumerGroupSession) MemberID() string {
	args := cgs.Called()

	return args.String(0)
}

func (cgs *ConsumerGroupSession) GenerationID() int32 {
	args := cgs.Called()

	return int32(args.Int(0))
}

func (cgs *ConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	cgs.Called(topic, partition, offset, metadata)
}

func (cgs *ConsumerGroupSession) Commit() {
	cgs.Called()
}

func (cgs *ConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	cgs.Called(topic, partition, offset, metadata)
}

func (cgs *ConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	cgs.Called(msg, metadata)
}

func (cgs *ConsumerGroupSession) Context() context.Context {
	args := cgs.Called()

	return args.Get(0).(context.Context)
}

type ConsumerGroupClaim struct {
	mock.Mock
}

func (cgc *ConsumerGroupClaim) Topic() string {
	args := cgc.Called()

	return args.String(0)
}

func (cgc *ConsumerGroupClaim) Partition() int32 {
	args := cgc.Called()

	return int32(args.Int(0))
}

func (cgc *ConsumerGroupClaim) InitialOffset() int64 {
	args := cgc.Called()

	return int64(args.Int(0))
}

func (cgc *ConsumerGroupClaim) HighWaterMarkOffset() int64 {
	args := cgc.Called()

	return int64(args.Int(0))
}

func (cgc *ConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	args := cgc.Called()

	out := make(chan *sarama.ConsumerMessage)

	go func() {
		if args.Get(0) == nil {
			out <- nil
		} else {
			out <- args.Get(0).(*sarama.ConsumerMessage)
		}

		close(out)
	}()

	return out
}
