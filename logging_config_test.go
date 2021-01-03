package config

import (
	"os"
	"testing"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggingConfigLoadLogs(t *testing.T) {
	os.Setenv("X", "foo")
	os.Setenv("Y", "123")
	os.Setenv("W", `["bar", "baz", "bonk"]`)

	type C struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"Q"`
	}

	logger := NewMockLogger()
	lc := NewConfig(NewEnvSourcer(""), WithLogger(logger))

	chunk := &C{}
	require.Nil(t, lc.Load(chunk))
	mockassert.CalledOnceWith(t, logger.PrintfFunc, mockassert.Values("Config loaded: %s", "\nQ=[\"bar\",\"baz\",\"bonk\"]\nX=foo\nY=123"))
}

func TestLoggingConfigMask(t *testing.T) {
	os.Setenv("X", "foo")
	os.Setenv("Y", "123")
	os.Setenv("W", `["bar", "baz", "bonk"]`)

	type C struct {
		X string   `env:"x"`
		Y int      `env:"y" mask:"true"`
		Z []string `env:"w" mask:"true"`
	}

	logger := NewMockLogger()
	lc := NewConfig(NewEnvSourcer(""), WithLogger(logger))

	chunk := &C{}
	require.Nil(t, lc.Load(chunk))
	mockassert.CalledOnceWith(t, logger.PrintfFunc, mockassert.Values("Config loaded: %s", "\nX=foo"))
}

func TestLoggingConfigBadMaskTag(t *testing.T) {
	type C struct {
		X string `env:"x" mask:"34"`
	}

	logger := NewMockLogger()
	lc := NewConfig(NewEnvSourcer(""), WithLogger(logger))

	chunk := &C{}
	assert.EqualError(t, lc.Load(chunk), "failed to serialize config (field 'X' has an invalid mask tag)")
}
