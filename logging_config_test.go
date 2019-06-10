package config

import (
	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"
)

type LoggingConfigSuite struct{}

func (s *LoggingConfigSuite) TestLoadLogs(t sweet.T) {
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

	Expect(lc.Load(chunk)).To(BeNil())
	Expect(logger.PrintfFunc).To(BeCalledOnceWith(
		"Config loaded: %s",
		"\nQ=[\"bar\",\"baz\",\"bonk\"]\nX=foo\nY=123",
	))
}

func (s *LoggingConfigSuite) TestMask(t sweet.T) {
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

	Expect(lc.Load(chunk)).To(BeNil())
	Expect(logger.PrintfFunc).To(BeCalledOnceWith(
		"Config loaded: %s",
		"\nX=foo",
	))
}

func (s *LoggingConfigSuite) TestBadMaskTag(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger, nil)
		chunk  = &TestBadMaskTagConfig{}
	)

	Expect(lc.Load(chunk)).To(MatchError("" +
		"failed to serialize config" +
		" (" +
		"field 'X' has an invalid mask tag" +
		")",
	))
}

func (s *LoggingConfigSuite) TestMustLoadLogs(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger, nil)
		chunk  = &TestSimpleConfig{}
	)

	lc.MustLoad(chunk)
	Expect(logger.PrintfFunc).To(BeCalledOnce())
}
