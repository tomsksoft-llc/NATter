package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	d := path.Join(path.Dir(b))
	p := fmt.Sprintf(
		"%s%cnatter%ccfg%ctest.toml",
		filepath.Dir(filepath.Dir(d)), os.PathSeparator, os.PathSeparator, os.PathSeparator,
	)

	if err := Load(p); err != nil {
		t.Skip("skipping configuration test, test.toml config reading error: ", err)
	}

	if !Bool("CONFIG_TEST_ENABLE") {
		t.Skip("skipping configuration test, disabled")
	}

	assert.Equal(t, "hello", String("CONFIG_TEST_STRING"))
	assert.Equal(t, true, Bool("CONFIG_TEST_BOOL"))
	assert.Equal(t, []string{"value1", "value2"}, StringSlice("CONFIG_TEST_STRING_SLICE"))
}
