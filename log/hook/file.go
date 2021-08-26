package hook

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type File struct {
	path string
	sync.Mutex
}

func NewFile(path string) *File {
	return &File{path: path}
}

func (f *File) Fire(entry *logrus.Entry) error {
	line, err := entry.Bytes()

	if err != nil {
		return err
	}

	f.Lock()
	defer f.Unlock()

	h, err := os.OpenFile(f.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)

	if err != nil {
		return err
	}

	defer h.Close()

	_, err = h.Write(line)

	return err
}

func (f *File) Levels() []logrus.Level {
	return logrus.AllLevels
}
