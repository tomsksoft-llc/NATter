package errtpl

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrConnect(t *testing.T) {
	assert.EqualError(t, ErrConnect(errors.New("error"), "world"), "unable connect to world: error")
}

func TestErrClose(t *testing.T) {
	assert.EqualError(t, ErrClose(errors.New("error"), "world"), "unable close correct world: error")
}

func TestErrMarshal(t *testing.T) {
	assert.EqualError(t, ErrMarshal(errors.New("error"), "world"), "unable marshal world: error")
}
