// +build !windows

package log

import (
	"net"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	hsyslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/stretchr/testify/assert"
)

func TestSetConfigWithSyslog(t *testing.T) {
	listener, err := net.Listen("tcp", "localhost:5000")

	assert.Nil(t, err)

	defer listener.Close()

	err = SetConfig(&Config{
		Level: "info",
		Syslog: ConfigSyslog{
			Enable: true,
			Host:   "localhost:5000",
		},
	})

	assert.Nil(t, err)

	assert.IsType(t, &hsyslog.SyslogHook{}, logrus.StandardLogger().Hooks[logrus.DebugLevel][0])

	std := logrus.StandardLogger()
	std.Hooks = make(map[logrus.Level][]logrus.Hook)

	logrus.SetOutput(os.Stderr)
}

func TestSetConfigWithSyslogOnConnectError(t *testing.T) {
	err := SetConfig(&Config{
		Level: "info",
		Syslog: ConfigSyslog{
			Enable: true,
			Host:   "incorrect-address",
		},
	})

	assert.Error(t, err)

	std := logrus.StandardLogger()
	std.Hooks = make(map[logrus.Level][]logrus.Hook)

	logrus.SetOutput(os.Stderr)
}
