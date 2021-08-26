package log

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
}

type Fields map[string]interface{}

type ConfigFile struct {
	Enable bool
	Path   string
}

type ConfigStream struct {
	Enable bool
}

type ConfigSyslog struct {
	Enable bool
	Host   string
	Name   string
}

type Config struct {
	Level  string
	Syslog ConfigSyslog
	Stream ConfigStream
	File   ConfigFile
}

func SetConfig(cfg *Config) error {
	lvl, err := logrus.ParseLevel(cfg.Level)

	if err != nil {
		return err
	}

	logrus.SetLevel(lvl)
	logrus.SetOutput(ioutil.Discard)

	if cfg.Syslog.Enable {
		if err = setupSyslog(cfg.Syslog.Host, cfg.Syslog.Name); err != nil {
			return err
		}
	}

	if cfg.File.Enable {
		if err = setupFile(cfg.File.Path); err != nil {
			return err
		}
	}

	if cfg.Stream.Enable {
		setupStream()
	}

	return nil
}

func FormatStruct(s interface{}) string {
	return fmt.Sprintf("%+v", s)
}

func WithFields(fields Fields) *Entry {
	return &Entry{logrus.WithFields(logrus.Fields(fields))}
}

func (e *Entry) WithFields(fields Fields) *Entry {
	e.Entry = e.Entry.WithFields(logrus.Fields(fields))

	return e
}

func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Error(args ...interface{}) {
	logrus.Error(args...)
}

func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}
