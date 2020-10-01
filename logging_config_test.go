package config

import (
	"testing"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggingConfigLoadLogs(t *testing.T) {
	config := NewMockConfig()
	logger := NewMockLogger()
	lc := NewLoggingConfig(config, logger, nil)
	chunk := &TestSimpleConfig{}

	config.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...TagModifier) error {
		target.(*TestSimpleConfig).X = "foo"
		target.(*TestSimpleConfig).Y = 123
		target.(*TestSimpleConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	})

	require.Nil(t, lc.Load(chunk))
	mockassert.CalledOnceWith(t, logger.PrintfFunc, mockassert.Values("Config loaded: %s", "\nQ=[\"bar\",\"baz\",\"bonk\"]\nX=foo\nY=123"))
}

func TestLoggingConfigMask(t *testing.T) {
	config := NewMockConfig()
	logger := NewMockLogger()
	lc := NewLoggingConfig(config, logger, nil)
	chunk := &TestMaskConfig{}

	config.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...TagModifier) error {
		target.(*TestMaskConfig).X = "foo"
		target.(*TestMaskConfig).Y = 123
		target.(*TestMaskConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	})

	require.Nil(t, lc.Load(chunk))
	mockassert.CalledOnceWith(t, logger.PrintfFunc, mockassert.Values("Config loaded: %s", "\nX=foo"))
}

func TestLoggingConfigBadMaskTag(t *testing.T) {
	config := NewMockConfig()
	logger := NewMockLogger()
	lc := NewLoggingConfig(config, logger, nil)
	chunk := &TestBadMaskTagConfig{}

	assert.EqualError(t, lc.Load(chunk), "failed to serialize config (field 'X' has an invalid mask tag)")
}

func TestLoggingConfigMustLoadLogs(t *testing.T) {
	config := NewMockConfig()
	logger := NewMockLogger()
	lc := NewLoggingConfig(config, logger, nil)
	chunk := &TestSimpleConfig{}

	lc.MustLoad(chunk)
	mockassert.CalledOnce(t, logger.PrintfFunc)
}
