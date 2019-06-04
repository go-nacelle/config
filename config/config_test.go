package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/go-nacelle/config/internal/fixtures"
	"github.com/go-nacelle/config/sourcer"
	"github.com/go-nacelle/config/tags"
)

type ConfigSuite struct{}

func (s *ConfigSuite) SetUpTest(t sweet.T) {
	os.Clearenv()
}

func (s *ConfigSuite) TestSimpleConfig(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &fixtures.TestSimpleConfig{}
	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
	Expect(chunk.Z).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *ConfigSuite) TestNestedJSONDeserialization(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app",
		map[string]string{
			"APP_P1": `{"v_int": 3, "v_float": 3.14, "v_bool": true}`,
			"APP_P2": `{"v_int": 5, "v_float": 6.28, "v_bool": false}`,
		}))

	chunk := &fixtures.TestEmbeddedJSONConfig{}
	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.P1).To(Equal(&fixtures.TestJSONPayload{V1: 3, V2: 3.14, V3: true}))
	Expect(chunk.P2).To(Equal(&fixtures.TestJSONPayload{V1: 5, V2: 6.28, V3: false}))
}

func (s *ConfigSuite) TestRequired(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &fixtures.TestRequiredConfig{}

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"no value supplied for field 'X'" +
		")",
	))
}

func (s *ConfigSuite) TestRequiredBadTag(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &fixtures.TestBadRequiredConfig{}

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"field 'X' has an invalid required tag" +
		")",
	))
}

func (s *ConfigSuite) TestDefault(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &fixtures.TestDefaultConfig{}

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *ConfigSuite) TestBadType(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "123", // silently converted to string
		"APP_Y": "foo",
		"APP_W": `bar`,
	}))

	chunk := &fixtures.TestSimpleConfig{}

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"value supplied for field 'Y' cannot be coerced into the expected type" +
		", " +
		"value supplied for field 'Z' cannot be coerced into the expected type" +
		")",
	))
}

func (s *ConfigSuite) TestBadDefaultType(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &fixtures.TestBadDefaultConfig{}

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"default value for field 'X' cannot be coerced into the expected type" +
		")",
	))
}

func (s *ConfigSuite) TestPostLoadConfig(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "3",
	}))

	chunk := &fixtures.TestPostLoadConfig{}
	Expect(config.Load(chunk)).To(BeNil())

	config = NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "-4",
	}))

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"X must be positive" +
		")",
	))
}

func (s *ConfigSuite) TestUnsettableFields(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &fixtures.TestUnsettableConfig{}

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"field 'x' can not be set" +
		")",
	))
}

func (s *ConfigSuite) TestLoad(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &fixtures.TestSimpleConfig{}

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
	Expect(chunk.Z).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *ConfigSuite) TestLoadIsomorphicType(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &fixtures.TestSimpleConfig{}

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal("foo"))
	Expect(chunk.Y).To(Equal(123))
	Expect(chunk.Z).To(Equal([]string{"bar", "baz", "bonk"}))
}

func (s *ConfigSuite) TestLoadPostLoadWithConversion(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_DURATION": "3",
	}))

	chunk := &fixtures.TestPostLoadConversion{}

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.Duration).To(Equal(time.Second * 3))
}

func (s *ConfigSuite) TestLoadPostLoadWithTags(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_FOO_DURATION": "3",
	}))

	chunk := &fixtures.TestPostLoadConversion{}

	Expect(config.Load(chunk, tags.NewEnvTagPrefixer("foo"))).To(BeNil())
	Expect(chunk.Duration).To(Equal(time.Second * 3))
}

func (s *ConfigSuite) TestBadConfigObjectTypes(t sweet.T) {
	Expect(NewConfig(NewFakeSourcer("app", nil)).Load(nil)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"configuration target is not a pointer to struct" +
		")",
	))

	Expect(NewConfig(NewFakeSourcer("app", nil)).Load("foo")).To(MatchError("" +
		"failed to load config" +
		" (" +
		"configuration target is not a pointer to struct" +
		")",
	))
}

func (s *ConfigSuite) TestEmbeddedConfig(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_A": "1",
		"APP_B": "2",
		"APP_C": "3",
		"APP_X": "4",
		"APP_Y": "5",
	}))

	chunk := &fixtures.TestParentConfig{}

	Expect(config.Load(chunk)).To(BeNil())
	Expect(chunk.X).To(Equal(4))
	Expect(chunk.Y).To(Equal(5))
	Expect(chunk.A).To(Equal(1))
	Expect(chunk.B).To(Equal(2))
	Expect(chunk.C).To(Equal(3))
}

func (s *ConfigSuite) TestEmbeddedConfigWithTags(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_FOO_A": "1",
		"APP_FOO_B": "2",
		"APP_FOO_C": "3",
		"APP_FOO_X": "4",
		"APP_FOO_Y": "5",
	}))

	chunk := &fixtures.TestParentConfig{}

	Expect(config.Load(chunk, tags.NewEnvTagPrefixer("foo"))).To(BeNil())
	Expect(chunk.X).To(Equal(4))
	Expect(chunk.Y).To(Equal(5))
	Expect(chunk.A).To(Equal(1))
	Expect(chunk.B).To(Equal(2))
	Expect(chunk.C).To(Equal(3))
}

func (s *ConfigSuite) TestEmbeddedConfigPostLoad(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_A": "1",
		"APP_B": "3",
		"APP_C": "2",
		"APP_X": "4",
		"APP_Y": "5",
	}))

	chunk := &fixtures.TestParentConfig{}

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"fields must be increasing" +
		")",
	))
}

func (s *ConfigSuite) TestBadEmbeddedObjectType(t sweet.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &fixtures.TestBadParentConfig{}

	Expect(config.Load(chunk)).To(MatchError("" +
		"failed to load config" +
		" (" +
		"invalid embedded type in configuration struct" +
		")",
	))
}

//
// Helpers

func NewFakeSourcer(prefix string, env map[string]string) sourcer.Sourcer {
	mock := NewMockSourcer()

	mock.GetFunc.SetDefaultHook(func(values []string) (string, sourcer.SourcerFlag, error) {
		if values[0] == "" {
			return "", sourcer.FlagSkip, nil
		}

		envvars := []string{
			strings.ToUpper(fmt.Sprintf("%s_%s", prefix, values[0])),
			strings.ToUpper(values[0]),
		}

		for _, envvar := range envvars {
			if val, ok := env[envvar]; ok {
				return val, sourcer.FlagFound, nil
			}
		}

		return "", sourcer.FlagMissing, nil
	})

	mock.TagsFunc.SetDefaultHook(func() []string {
		return []string{"env"}
	})

	return mock
}
