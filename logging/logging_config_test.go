package logging

import (
	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"

	"github.com/efritz/zubrin/internal/fixtures"
	tags "github.com/efritz/zubrin/tags"
)

type LoggingConfigSuite struct{}

func (s *LoggingConfigSuite) TestLoadLogs(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &fixtures.TestSimpleConfig{}
	)

	config.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...tags.TagModifier) error {
		target.(*fixtures.TestSimpleConfig).X = "foo"
		target.(*fixtures.TestSimpleConfig).Y = 123
		target.(*fixtures.TestSimpleConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	})

	Expect(lc.Load(chunk)).To(BeNil())
	Expect(logger.PrintfFunc).To(BeCalledOnceWith(
		"Config loaded: %s",
		`Q=["bar","baz","bonk"], X=foo, Y=123`,
	))
}

func (s *LoggingConfigSuite) TestMask(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &fixtures.TestMaskConfig{}
	)

	config.LoadFunc.SetDefaultHook(func(target interface{}, modifiers ...tags.TagModifier) error {
		target.(*fixtures.TestMaskConfig).X = "foo"
		target.(*fixtures.TestMaskConfig).Y = 123
		target.(*fixtures.TestMaskConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	})

	Expect(lc.Load(chunk)).To(BeNil())
	Expect(logger.PrintfFunc).To(BeCalledOnceWith(
		"Config loaded: %s",
		"X=foo",
	))
}

func (s *LoggingConfigSuite) TestBadMaskTag(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &fixtures.TestBadMaskTagConfig{}
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
		lc     = NewLoggingConfig(config, logger)
		chunk  = &fixtures.TestSimpleConfig{}
	)

	lc.MustLoad(chunk)
	Expect(logger.PrintfFunc).To(BeCalledOnce())
}
