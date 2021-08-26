package errtpl

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")
)

func ErrConnect(err error, service string) error {
	return errors.Wrapf(err, "unable connect to %s", service)
}

func ErrClose(err error, service string) error {
	return errors.Wrapf(err, "unable close correct %s", service)
}

func ErrMarshal(err error, v interface{}) error {
	return errors.Wrapf(err, "unable marshal %+v", v)
}
