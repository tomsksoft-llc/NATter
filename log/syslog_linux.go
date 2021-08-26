// +build !windows

package log

import (
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	hsyslog "github.com/sirupsen/logrus/hooks/syslog"
)

func setupSyslog(host string, name string) error {
	hostname, err := os.Hostname()

	if err != nil {
		return err
	}

	tag := name + "#" + hostname

	hook, err := hsyslog.NewSyslogHook("udp", host, syslog.LOG_LOCAL0|syslog.LOG_DEBUG, tag)

	if err != nil {
		return err
	}

	logrus.AddHook(hook)

	return nil
}
