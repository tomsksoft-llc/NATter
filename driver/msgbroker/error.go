package msgbroker

import (
	"github.com/pkg/errors"
)

var (
	ErrNoResponders = errors.New("no responders available for request")
)

func ErrBadReply(err error, topic string) error {
	return errors.Wrapf(err, "bad reply from topic %s", topic)
}

func ErrPublish(err error, topic string) error {
	return errors.Wrapf(err, "unable publish to topic %s", topic)
}

func ErrSubscribe(err error, topic string) error {
	return errors.Wrapf(err, "unable subscribe to topic %s", topic)
}

func ErrUnsubscribe(err error, topic string) error {
	return errors.Wrapf(err, "unable unsubscribe from topic %s", topic)
}

func ErrRespond(err error, topic string) error {
	return errors.Wrapf(err, "unable respond to topic %s", topic)
}
