package config

import (
	"testing"

	mockassert "github.com/efritz/go-mockgen/assert"
	"github.com/stretchr/testify/assert"
)

func TestLoggingConfigLoadLogs(t *testing.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger, nil)
		chunk  = &TestSimpleConfig{}
	)

	config.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...TagModifier) error {
		target.(*TestSimpleConfig).X = "foo"
		target.(*TestSimpleConfig).Y = 123
		target.(*TestSimpleConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	})

	assert.Nil(t, lc.Load(chunk))
	mockassert.CalledOnceMatching(t, logger.PrintfFunc, func(t assert.TestingT, call interface{}) bool {
		c := call.(LoggerPrintfFuncCall)
		// TODO - also match "\nQ=[\"bar\",\"baz\",\"bonk\"]\nX=foo\nY=123"
		return c.Arg0 == "Config loaded: %s"
	})
}

func TestLoggingConfigMask(t *testing.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger, nil)
		chunk  = &TestMaskConfig{}
	)

	config.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...TagModifier) error {
		target.(*TestMaskConfig).X = "foo"
		target.(*TestMaskConfig).Y = 123
		target.(*TestMaskConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	})

	assert.Nil(t, lc.Load(chunk))
	mockassert.CalledOnceMatching(t, logger.PrintfFunc, func(t assert.TestingT, call interface{}) bool {
		c := call.(LoggerPrintfFuncCall)
		// TODO - also match "\nX=foo"
		return c.Arg0 == "Config loaded: %s"
	})
}

func TestLoggingConfigBadMaskTag(t *testing.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger, nil)
		chunk  = &TestBadMaskTagConfig{}
	)

	assert.EqualError(t, lc.Load(chunk), "failed to serialize config (field 'X' has an invalid mask tag)")
}

func TestLoggingConfigMustLoadLogs(t *testing.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger, nil)
		chunk  = &TestSimpleConfig{}
	)

	lc.MustLoad(chunk)
	mockassert.CalledOnce(t, logger.PrintfFunc)
}
