package zubrin

//go:generate go-mockgen github.com/efritz/zubrin -i Config -o mock_config_test.go -f
//go:generate go-mockgen github.com/efritz/zubrin -i Logger -o mock_logger_test.go -f

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type LoggingConfigSuite struct{}

func (s *LoggingConfigSuite) TestLoadLogs(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &TestSimpleConfig{}
	)

	config.LoadFunc = func(target interface{}, modifiers ...TagModifier) error {
		target.(*TestSimpleConfig).X = "foo"
		target.(*TestSimpleConfig).Y = 123
		target.(*TestSimpleConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	}

	Expect(lc.Load(chunk)).To(BeNil())
	Expect(logger.PrintfFuncCallCount()).To(Equal(1))

	params := logger.PrintfFuncCallParams()[0]
	Expect(params.Arg0).To(Equal("Config loaded from environment: %s"))
	Expect(params.Arg1[0]).To(Equal(`Q=["bar","baz","bonk"], X=foo, Y=123`))
}

func (s *LoggingConfigSuite) TestMask(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
		chunk  = &TestMaskConfig{}
	)

	config.LoadFunc = func(target interface{}, modifiers ...TagModifier) error {
		target.(*TestMaskConfig).X = "foo"
		target.(*TestMaskConfig).Y = 123
		target.(*TestMaskConfig).Z = []string{"bar", "baz", "bonk"}
		return nil
	}

	Expect(lc.Load(chunk)).To(BeNil())
	Expect(logger.PrintfFuncCallCount()).To(Equal(1))

	params := logger.PrintfFuncCallParams()[0]
	Expect(params.Arg0).To(Equal("Config loaded from environment: %s"))
	Expect(params.Arg1[0]).To(Equal("X=foo"))
}

func (s *LoggingConfigSuite) TestBadMaskTag(t sweet.T) {
	var (
		config = NewMockConfig()
		logger = NewMockLogger()
		lc     = NewLoggingConfig(config, logger)
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
		lc     = NewLoggingConfig(config, logger)
		chunk  = &TestSimpleConfig{}
	)

	lc.MustLoad(chunk)
	Expect(logger.PrintfFuncCallCount()).To(Equal(1))
}
