package config

import (
	"testing"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggingConfigLoadLogs(t *testing.T) {
	type C struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"Q"`
	}

	config := NewMockConfig()
	logger := NewMockLogger()
	lc := NewLoggingConfig(config, logger, nil)
	chunk := &C{}

	config.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...TagModifier) error {
		target.(*C).X = "foo"
		target.(*C).Y = 123
		target.(*C).Z = []string{"bar", "baz", "bonk"}
		return nil
	})

	require.Nil(t, lc.Load(chunk))
	mockassert.CalledOnceWith(t, logger.PrintfFunc, mockassert.Values("Config loaded: %s", "\nQ=[\"bar\",\"baz\",\"bonk\"]\nX=foo\nY=123"))
}

func TestLoggingConfigMask(t *testing.T) {
	type C struct {
		X string   `env:"x"`
		Y int      `env:"y" mask:"true"`
		Z []string `env:"w" mask:"true"`
	}

	config := NewMockConfig()
	logger := NewMockLogger()
	lc := NewLoggingConfig(config, logger, nil)
	chunk := &C{}

	config.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...TagModifier) error {
		target.(*C).X = "foo"
		target.(*C).Y = 123
		target.(*C).Z = []string{"bar", "baz", "bonk"}
		return nil
	})

	require.Nil(t, lc.Load(chunk))
	mockassert.CalledOnceWith(t, logger.PrintfFunc, mockassert.Values("Config loaded: %s", "\nX=foo"))
}

func TestLoggingConfigBadMaskTag(t *testing.T) {
	type C struct {
		X string `env:"x" mask:"34"`
	}

	config := NewMockConfig()
	logger := NewMockLogger()
	lc := NewLoggingConfig(config, logger, nil)
	chunk := &C{}

	assert.EqualError(t, lc.Load(chunk), "failed to serialize config (field 'X' has an invalid mask tag)")
}

func TestLoggingConfigMustLoadLogs(t *testing.T) {
	type C struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"Q"`
	}

	config := NewMockConfig()
	logger := NewMockLogger()
	lc := NewLoggingConfig(config, logger, nil)
	chunk := &C{}

	lc.MustLoad(chunk)
	mockassert.CalledOnce(t, logger.PrintfFunc)
}
