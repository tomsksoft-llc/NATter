package log

import (
	"io/ioutil"
	"os"
	"testing"

	"NATter/log/hook"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/sirupsen/logrus/hooks/writer"
	"github.com/stretchr/testify/assert"
)

func TestSetConfig(t *testing.T) {
	err := SetConfig(&Config{
		Level: "info",
	})

	assert.Nil(t, err)
	assert.Equal(t, ioutil.Discard, logrus.StandardLogger().Out)

	logrus.SetOutput(os.Stderr)
}

func TestSetConfigOnEnableStream(t *testing.T) {
	err := SetConfig(&Config{
		Level: "info",
		Stream: ConfigStream{
			Enable: true,
		},
	})

	assert.Nil(t, err)
	assert.IsType(t, &writer.Hook{}, logrus.StandardLogger().Hooks[logrus.InfoLevel][0])
	assert.IsType(t, &writer.Hook{}, logrus.StandardLogger().Hooks[logrus.PanicLevel][0])

	std := logrus.StandardLogger()
	std.Hooks = make(map[logrus.Level][]logrus.Hook)

	logrus.SetOutput(os.Stderr)
}

func TestSetConfigOnEnableFile(t *testing.T) {
	err := SetConfig(&Config{
		Level: "info",
		File: ConfigFile{
			Enable: true,
			Path:   "/logs",
		},
	})

	assert.Nil(t, err)
	assert.IsType(t, &hook.File{}, logrus.StandardLogger().Hooks[logrus.InfoLevel][0])

	std := logrus.StandardLogger()
	std.Hooks = make(map[logrus.Level][]logrus.Hook)

	logrus.SetOutput(os.Stderr)
}

func TestSetConfigOnEnableFileOnNoPath(t *testing.T) {
	err := SetConfig(&Config{
		Level: "info",
		File: ConfigFile{
			Enable: true,
		},
	})

	assert.Error(t, err)

	std := logrus.StandardLogger()
	std.Hooks = make(map[logrus.Level][]logrus.Hook)

	logrus.SetOutput(os.Stderr)
}

func TestSetConfigOnEnableFileOnNotWritable(t *testing.T) {
	err := SetConfig(&Config{
		Level: "info",
		File: ConfigFile{
			Enable: true,
			Path:   "/",
		},
	})

	assert.Error(t, err)

	std := logrus.StandardLogger()
	std.Hooks = make(map[logrus.Level][]logrus.Hook)

	logrus.SetOutput(os.Stderr)
}

func TestSetConfigOnIncorrectLevel(t *testing.T) {
	err := SetConfig(&Config{
		Level: "unknown_level",
	})

	assert.NotNil(t, err)
}

func TestFormatStruct(t *testing.T) {
	a := struct{ hello string }{hello: "world"}
	assert.Equal(t, "{hello:world}", FormatStruct(a))
}

func TestEntryWithFields(t *testing.T) {
	ent := WithFields(Fields{
		"test": true,
	}).WithFields(Fields{
		"test1": true,
	})

	assert.Equal(t, logrus.Fields{
		"test":  true,
		"test1": true,
	}, ent.Data)
}

func TestEntryTrace(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Trace("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.TraceLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryTracef(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Tracef("hello %s", "!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.TraceLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello !", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryDebug(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Debug("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryDebugf(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Debugf("hello %s", "!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello !", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryInfo(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Info("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryInfof(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Infof("hello %s", "!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello !", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryWarn(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Warn("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryWarnf(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Warnf("hello %s", "!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello !", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryError(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Error("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryErrorf(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.TraceLevel)

	WithFields(Fields{
		"test": true,
	}).Errorf("hello %s", "!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello !", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryFatal(t *testing.T) {
	hook := test.NewGlobal() //nolint:staticcheck //false positive

	logrus.SetLevel(logrus.FatalLevel)

	logrus.StandardLogger().ExitFunc = func(int) {}

	WithFields(Fields{
		"test": true,
	}).Fatal("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.FatalLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryFatalf(t *testing.T) {
	hook := test.NewGlobal() //nolint:staticcheck //false positive

	logrus.SetLevel(logrus.FatalLevel)

	logrus.StandardLogger().ExitFunc = func(int) {}

	WithFields(Fields{
		"test": true,
	}).Fatalf("hello %s", "!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.FatalLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello !", hook.LastEntry().Message)
	assert.Equal(t, logrus.Fields{
		"test": true,
	}, hook.LastEntry().Data)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestEntryPanic(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.PanicLevel)

	defer func() {
		err := recover()
		assert.NotNil(t, err)

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.PanicLevel, hook.LastEntry().Level)
		assert.Equal(t, "hello", hook.LastEntry().Message)
		assert.Equal(t, logrus.Fields{
			"test": true,
		}, hook.LastEntry().Data)

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	}()

	WithFields(Fields{
		"test": true,
	}).Panic("hello")
}

func TestEntryPanicf(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.PanicLevel)

	defer func() {
		err := recover()
		assert.NotNil(t, err)

		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.PanicLevel, hook.LastEntry().Level)
		assert.Equal(t, "hello !", hook.LastEntry().Message)
		assert.Equal(t, logrus.Fields{
			"test": true,
		}, hook.LastEntry().Data)

		hook.Reset()
		assert.Nil(t, hook.LastEntry())
	}()

	WithFields(Fields{
		"test": true,
	}).Panicf("hello %s", "!")
}

func TestDebug(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.DebugLevel)

	Debug("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestDebugf(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.DebugLevel)

	Debugf("hello %s", "!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello !", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestInfo(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.InfoLevel)

	Info("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestWarnf(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.WarnLevel)

	Warnf("hello %s", "!")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello !", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestError(t *testing.T) {
	hook := test.NewGlobal()

	logrus.SetLevel(logrus.ErrorLevel)

	Error("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}

func TestFatal(t *testing.T) {
	hook := test.NewGlobal() //nolint:staticcheck //false positive

	logrus.SetLevel(logrus.FatalLevel)

	logrus.StandardLogger().ExitFunc = func(int) {}

	Fatal("hello")

	assert.Equal(t, 1, len(hook.Entries))
	assert.Equal(t, logrus.FatalLevel, hook.LastEntry().Level)
	assert.Equal(t, "hello", hook.LastEntry().Message)

	hook.Reset()
	assert.Nil(t, hook.LastEntry())
}
