package mock

import (
	"github.com/pkg/errors"
)

type ReadCloserErr struct {
}

func (er *ReadCloserErr) Read([]byte) (int, error) {
	return 0, errors.New("error")
}

func (er *ReadCloserErr) Close() error {
	return errors.New("error")
}
