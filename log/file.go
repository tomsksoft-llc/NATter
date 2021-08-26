package log

import (
	"os"

	"NATter/log/hook"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func setupFile(path string) error {
	if len(path) == 0 {
		return errors.New("logger path couldn't be empty")
	}

	h, err := os.OpenFile(path, os.O_WRONLY, 0664)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	defer h.Close()

	logrus.AddHook(hook.NewFile(path))

	return nil
}
