package msgbroker

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestLogErrorHandle(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.ErrorLevel)

	err := errors.New("error")
	LogErrorHandle(err, "topic")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "unable handle message", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"topic": "topic",
		"error": err,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())

	logrus.SetLevel(logrus.InfoLevel)
}

func TestLogDebugSubscribed(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.DebugLevel)

	LogDebugSubscribed("topic")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "subscribed", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"topic": "topic",
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())

	logrus.SetLevel(logrus.InfoLevel)
}

func TestLogDebugUnsubscribed(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.DebugLevel)

	LogDebugUnsubscribed("topic")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "unsubscribed", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"topic": "topic",
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())

	logrus.SetLevel(logrus.InfoLevel)
}

func TestLogDebugPublished(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.DebugLevel)

	LogDebugPublished("topic", "payload")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "published", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"payload": "payload",
		"topic":   "topic",
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())

	logrus.SetLevel(logrus.InfoLevel)
}

func TestLogDebugRequested(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.DebugLevel)

	LogDebugRequested("topic", "req", "resp")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "requested", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"request":  "req",
		"response": "resp",
		"topic":    "topic",
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())

	logrus.SetLevel(logrus.InfoLevel)
}

func TestLogDebugReceived(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.DebugLevel)

	LogDebugReceived("topic", "payload")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "received", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"payload": "payload",
		"topic":   "topic",
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())

	logrus.SetLevel(logrus.InfoLevel)
}

func TestLogDebugResponded(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.DebugLevel)

	LogDebugResponded("topic", "payload")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "responded", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"payload": "payload",
		"topic":   "topic",
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())

	logrus.SetLevel(logrus.InfoLevel)
}
