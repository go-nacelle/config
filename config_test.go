package config

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigSimpleConfig(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &TestSimpleConfig{}
	assert.Nil(t, config.Load(chunk))
	assert.Equal(t, "foo", chunk.X)
	assert.Equal(t, 123, chunk.Y)
	assert.Equal(t, []string{"bar", "baz", "bonk"}, chunk.Z)
}

func TestConfigNestedJSONDeserialization(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app",
		map[string]string{
			"APP_P1": `{"v_int": 3, "v_float": 3.14, "v_bool": true}`,
			"APP_P2": `{"v_int": 5, "v_float": 6.28, "v_bool": false}`,
		}))

	chunk := &TestEmbeddedJSONConfig{}
	assert.Nil(t, config.Load(chunk))
	assert.Equal(t, &TestJSONPayload{V1: 3, V2: 3.14, V3: true}, chunk.P1)
	assert.Equal(t, &TestJSONPayload{V1: 5, V2: 6.28, V3: false}, chunk.P2)
}

func TestConfigRequired(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &TestRequiredConfig{}
	assert.EqualError(t, config.Load(chunk), "failed to load config (no value supplied for field 'X')")
}

func TestConfigRequiredBadTag(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &TestBadRequiredConfig{}
	assert.EqualError(t, config.Load(chunk), "failed to load config (field 'X' has an invalid required tag)")
}

func TestConfigDefault(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &TestDefaultConfig{}

	assert.Nil(t, config.Load(chunk))
	assert.Equal(t, "foo", chunk.X)
	assert.Equal(t, []string{"bar", "baz", "bonk"}, chunk.Y)
}

func TestConfigBadType(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "123", // silently converted to string
		"APP_Y": "foo",
		"APP_W": `bar`,
	}))

	chunk := &TestSimpleConfig{}
	assert.EqualError(t, config.Load(chunk), fmt.Sprintf("failed to load config (%s)", strings.Join([]string{
		"value supplied for field 'Y' cannot be coerced into the expected type",
		"value supplied for field 'Z' cannot be coerced into the expected type",
	}, ", ")))
}

func TestConfigBadDefaultType(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &TestBadDefaultConfig{}
	assert.EqualError(t, config.Load(chunk), "failed to load config (default value for field 'X' cannot be coerced into the expected type)")
}

func TestConfigPostLoadConfig(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "3",
	}))

	chunk := &TestPostLoadConfig{}
	assert.Nil(t, config.Load(chunk))

	config = NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "-4",
	}))

	assert.EqualError(t, config.Load(chunk), "failed to load config (X must be positive)")
}

func TestConfigUnsettableFields(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &TestUnsettableConfig{}
	assert.EqualError(t, config.Load(chunk), "failed to load config (field 'x' can not be set)")
}

func TestConfigLoad(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &TestSimpleConfig{}

	assert.Nil(t, config.Load(chunk))
	assert.Equal(t, "foo", chunk.X)
	assert.Equal(t, 123, chunk.Y)
	assert.Equal(t, []string{"bar", "baz", "bonk"}, chunk.Z)
}

func TestConfigLoadIsomorphicType(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &TestSimpleConfig{}

	assert.Nil(t, config.Load(chunk))
	assert.Equal(t, "foo", chunk.X)
	assert.Equal(t, 123, chunk.Y)
	assert.Equal(t, []string{"bar", "baz", "bonk"}, chunk.Z)
}

func TestConfigLoadPostLoadWithConversion(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_DURATION": "3",
	}))

	chunk := &TestPostLoadConversion{}

	assert.Nil(t, config.Load(chunk))
	assert.Equal(t, time.Second*3, chunk.Duration)
}

func TestConfigLoadPostLoadWithTags(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_FOO_DURATION": "3",
	}))

	chunk := &TestPostLoadConversion{}

	assert.Nil(t, config.Load(chunk, NewEnvTagPrefixer("foo")))
	assert.Equal(t, time.Second*3, chunk.Duration)
}

func TestConfigBadConfigObjectTypes(t *testing.T) {
	assert.EqualError(t, NewConfig(NewFakeSourcer("app", nil)).Load(nil), "failed to load config (configuration target is not a pointer to struct)")
	assert.EqualError(t, NewConfig(NewFakeSourcer("app", nil)).Load("foo"), "failed to load config (configuration target is not a pointer to struct)")
}

func TestConfigEmbeddedConfig(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_A": "1",
		"APP_B": "2",
		"APP_C": "3",
		"APP_X": "4",
		"APP_Y": "5",
	}))

	chunk := &TestParentConfig{}

	assert.Nil(t, config.Load(chunk))
	assert.Equal(t, 4, chunk.X)
	assert.Equal(t, 5, chunk.Y)
	assert.Equal(t, 1, chunk.A)
	assert.Equal(t, 2, chunk.B)
	assert.Equal(t, 3, chunk.C)
}

func TestConfigEmbeddedConfigWithTags(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_FOO_A": "1",
		"APP_FOO_B": "2",
		"APP_FOO_C": "3",
		"APP_FOO_X": "4",
		"APP_FOO_Y": "5",
	}))

	chunk := &TestParentConfig{}

	assert.Nil(t, config.Load(chunk, NewEnvTagPrefixer("foo")))
	assert.Equal(t, 4, chunk.X)
	assert.Equal(t, 5, chunk.Y)
	assert.Equal(t, 1, chunk.A)
	assert.Equal(t, 2, chunk.B)
	assert.Equal(t, 3, chunk.C)
}

func TestConfigEmbeddedConfigPostLoad(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_A": "1",
		"APP_B": "3",
		"APP_C": "2",
		"APP_X": "4",
		"APP_Y": "5",
	}))

	chunk := &TestParentConfig{}
	assert.EqualError(t, config.Load(chunk), "failed to load config (fields must be increasing)")
}

func TestConfigBadEmbeddedObjectType(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &TestBadParentConfig{}
	assert.EqualError(t, config.Load(chunk), "failed to load config (invalid embedded type in configuration struct)")
}

//
// Helpers

func NewFakeSourcer(prefix string, env map[string]string) Sourcer {
	mock := NewMockSourcer()

	mock.GetFunc.SetDefaultHook(func(values []string) (string, SourcerFlag, error) {
		if values[0] == "" {
			return "", FlagSkip, nil
		}

		envvars := []string{
			strings.ToUpper(fmt.Sprintf("%s_%s", prefix, values[0])),
			strings.ToUpper(values[0]),
		}

		for _, envvar := range envvars {
			if val, ok := env[envvar]; ok {
				return val, FlagFound, nil
			}
		}

		return "", FlagMissing, nil
	})

	mock.TagsFunc.SetDefaultHook(func() []string {
		return []string{"env"}
	})

	return mock
}
