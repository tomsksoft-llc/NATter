// +build windows

package log

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetConfigOnEnableSyslog(t *testing.T) {
	err := SetConfig(&Config{
		Level: "info",
		Syslog: ConfigSyslog{
			Enable: true,
			Host:   "localhost:5000",
		},
	})
	assert.Nil(t, err)

	logrus.SetOutput(os.Stderr)
}
