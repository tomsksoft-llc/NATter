package protobuf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func TestEncoderMarshal(t *testing.T) {
	enc := &Encoder{}

	msg1, err := proto.Marshal(&Batch{
		Messages: []*anypb.Any{{}},
	})

	assert.Nil(t, err)

	msg2, err := proto.Marshal(&Batch{
		Messages: []*anypb.Any{{}, {}},
	})

	assert.Nil(t, err)

	batch, err := proto.Marshal(&Batch{
		Messages: []*anypb.Any{{}, {}, {}},
	})

	assert.Nil(t, err)

	res, err := enc.Marshal([][]byte{
		msg1,
		msg2,
	})

	assert.Equal(t, batch, res)
	assert.Nil(t, err)
}

func TestEncoderMarshalOnUmarshalError(t *testing.T) {
	enc := &Encoder{}

	res, err := enc.Marshal([][]byte{
		[]byte("broken-data"),
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}
