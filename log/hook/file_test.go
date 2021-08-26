package hook

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFire(t *testing.T) {
	dir, err := ioutil.TempDir(".", "")
	assert.Nil(t, err)

	defer os.RemoveAll(dir)

	hook := NewFile(dir + string(os.PathSeparator) + "test.log")

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.AddHook(hook)

	logger.Info("hello")

	data, err := ioutil.ReadFile(dir + string(os.PathSeparator) + "test.log")
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(data), "hello"))
}

func TestFireOnError(t *testing.T) {
	hook := NewFile(string(os.PathSeparator) + "thisdirrectorydoesntexist" + string(os.PathSeparator))

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.AddHook(hook)

	logger.Info("hello")
}
