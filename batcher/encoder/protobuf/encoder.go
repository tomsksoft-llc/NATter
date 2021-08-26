package protobuf

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type Encoder struct {
}

func (e *Encoder) Marshal(msgs [][]byte) ([]byte, error) {
	accumulator := &Batch{
		Messages: []*anypb.Any{},
	}

	for _, msg := range msgs {
		batch := &Batch{}

		if err := proto.Unmarshal(msg, batch); err != nil {
			return nil, err
		}

		accumulator.Messages = append(accumulator.Messages, batch.Messages...)
	}

	b, err := proto.Marshal(accumulator)

	if err != nil {
		return nil, err
	}

	return b, nil
}
