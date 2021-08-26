package kafka

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"NATter/config"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	p := fmt.Sprintf("%s%ccfg%ctest.toml", filepath.Dir(filepath.Dir(filepath.Dir(d))), os.PathSeparator, os.PathSeparator)

	if err := config.Load(p); err != nil {
		t.Skip("skipping kafka connection test, test.toml config reading error: ", err)
	}

	if !config.Bool("KAFKA_TEST_ENABLE") {
		t.Skip("kafka connection test, disabled")
	}

	conn, err := NewConn(&ConnConfig{
		Version: config.String("KAFKA_VERSION"),
		Servers: config.StringSlice("KAFKA_SERVERS"),
		Group:   "test_group",
	})

	assert.NotNil(t, conn)
	assert.Nil(t, err)
}
