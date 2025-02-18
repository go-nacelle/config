package config

import (
	"fmt"
	"strings"
	"testing"
	"time"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigSimpleConfig(t *testing.T) {
	type C struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"Q"`
	}

	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &C{}
	require.Nil(t, config.Load(chunk))
	assert.Equal(t, "foo", chunk.X)
	assert.Equal(t, 123, chunk.Y)
	assert.Equal(t, []string{"bar", "baz", "bonk"}, chunk.Z)
}

func TestConfigNestedJSONDeserialization(t *testing.T) {
	type P struct {
		V1 int     `json:"v_int"`
		V2 float64 `json:"v_float"`
		V3 bool    `json:"v_bool"`
	}
	type C struct {
		P1 *P `env:"p1"`
		P2 *P `env:"p2"`
	}

	config := NewConfig(NewFakeSourcer("app",
		map[string]string{
			"APP_P1": `{"v_int": 3, "v_float": 3.14, "v_bool": true}`,
			"APP_P2": `{"v_int": 5, "v_float": 6.28, "v_bool": false}`,
		}))

	chunk := &C{}
	require.Nil(t, config.Load(chunk))
	assert.Equal(t, &P{V1: 3, V2: 3.14, V3: true}, chunk.P1)
	assert.Equal(t, &P{V1: 5, V2: 6.28, V3: false}, chunk.P2)
}

func TestConfigRequired(t *testing.T) {
	type C struct {
		X string `env:"x" required:"true"`
	}

	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &C{}
	assert.EqualError(t,
		config.Load(chunk),
		"failed to load config: no value supplied for field 'X'",
	)
}

func TestConfigRequiredBadTag(t *testing.T) {
	type C struct {
		X string `env:"x" required:"yup"`
	}

	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &C{}
	assert.EqualError(t,
		config.Load(chunk),
		"failed to load config: field 'X' has an invalid required tag",
	)
}

func TestConfigDefault(t *testing.T) {
	type C struct {
		X string   `env:"x" default:"foo"`
		Y []string `env:"y" default:"[\"bar\", \"baz\", \"bonk\"]"`
	}

	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &C{}

	require.Nil(t, config.Load(chunk))
	assert.Equal(t, "foo", chunk.X)
	assert.Equal(t, []string{"bar", "baz", "bonk"}, chunk.Y)
}

func TestConfigBadType(t *testing.T) {
	type C struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"Q"`
	}

	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "123", // silently converted to string
		"APP_Y": "foo",
		"APP_W": `bar`,
	}))

	chunk := &C{}

	var errLoad *LoadError
	require.ErrorAs(t, config.Load(chunk), &errLoad)
	assert.ElementsMatch(t,
		errLoad.Errors,
		[]error{
			fmt.Errorf("value supplied for field 'Y' cannot be coerced into the expected type"),
			fmt.Errorf("value supplied for field 'Z' cannot be coerced into the expected type"),
		},
	)
}

func TestConfigBadDefaultType(t *testing.T) {
	type C struct {
		X int `env:"x" default:"foo"`
	}

	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &C{}

	var errLoad *LoadError
	require.ErrorAs(t, config.Load(chunk), &errLoad)
	assert.ElementsMatch(t,
		errLoad.Errors,
		[]error{
			fmt.Errorf("default value for field 'X' cannot be coerced into the expected type"),
		},
	)
}

func TestConfigPostLoadConfig(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "3",
	}))

	chunk := &testPostLoadConfig{}
	require.Nil(t, config.Load(chunk))

	config = NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "-4",
	}))

	var postErr *PostLoadError
	require.ErrorAs(t, config.Load(chunk), &postErr)
	assert.EqualError(t,
		postErr.Unwrap(),
		"X must be positive",
	)
}

type testPostLoadConfig struct {
	X int `env:"X"`
}

func (c *testPostLoadConfig) PostLoad() error {
	if c.X < 0 {
		return fmt.Errorf("X must be positive")
	}

	return nil
}

func TestConfigUnsettableFields(t *testing.T) {
	type C struct {
		x int `env:"s"`
	}

	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &C{}

	var errLoad *LoadError
	require.ErrorAs(t, config.Load(chunk), &errLoad)
	assert.ElementsMatch(t,
		errLoad.Errors,
		[]error{
			fmt.Errorf("field 'x' can not be set"),
		},
	)
}

func TestConfigLoad(t *testing.T) {
	type C struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"Q"`
	}

	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &C{}

	require.Nil(t, config.Load(chunk))
	assert.Equal(t, "foo", chunk.X)
	assert.Equal(t, 123, chunk.Y)
	assert.Equal(t, []string{"bar", "baz", "bonk"}, chunk.Z)
}

func TestConfigLoadDumpOption(t *testing.T) {
	t.Run("defaults to false", func(t *testing.T) {
		logger := NewMockLogger()

		c := NewConfig(
			NewFakeSourcer("app", nil),
			WithLogger(logger),
		)

		assert.Error(t, c.Load(nil))

		mockassert.NotCalled(t, logger.PrintfFunc)
	})
	t.Run("can be turned on", func(t *testing.T) {
		logger := NewMockLogger()

		c := NewConfig(
			NewFakeSourcer("app", nil),
			WithSourceDumpOnError(true),
			WithLogger(logger),
		)

		assert.Error(t, c.Load(nil))

		assert.Equal(t,
			[]LoggerPrintfFuncCall{
				{Arg0: "Config source assets: %s", Arg1: []interface{}{""}},
				{Arg0: "Config source contents: %s", Arg1: []interface{}{"<no values>"}},
			},
			logger.PrintfFunc.History(),
		)
	})
}

func TestConfigLoadIsomorphicType(t *testing.T) {
	type C struct {
		X string   `env:"x"`
		Y int      `env:"y"`
		Z []string `env:"w" display:"Q"`
	}

	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_X": "foo",
		"APP_Y": "123",
		"APP_W": `["bar", "baz", "bonk"]`,
	}))

	chunk := &C{}

	require.Nil(t, config.Load(chunk))
	assert.Equal(t, "foo", chunk.X)
	assert.Equal(t, 123, chunk.Y)
	assert.Equal(t, []string{"bar", "baz", "bonk"}, chunk.Z)
}

func TestConfigLoadPostLoadWithConversion(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_DURATION": "3",
	}))

	chunk := &testPostLoadConversion{}
	require.Nil(t, config.Load(chunk))
	assert.Equal(t, time.Second*3, chunk.Duration)
}

func TestConfigLoadPostLoadWithTags(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_FOO_DURATION": "3",
	}))

	chunk := &testPostLoadConversion{}
	require.Nil(t, config.Load(chunk, NewEnvTagPrefixer("foo")))
	assert.Equal(t, time.Second*3, chunk.Duration)
}

type testPostLoadConversion struct {
	RawDuration int `env:"duration"`
	Duration    time.Duration
}

func (c *testPostLoadConversion) PostLoad() error {
	c.Duration = time.Duration(c.RawDuration) * time.Second
	return nil
}

func TestConfigBadConfigObjectTypes(t *testing.T) {
	t.Run("nil config object", func(t *testing.T) {
		var errLoad *LoadError
		require.ErrorAs(t,
			NewConfig(NewFakeSourcer("app", nil)).Load(nil),
			&errLoad,
		)
		assert.ElementsMatch(t,
			errLoad.Errors,
			[]error{
				fmt.Errorf("configuration target is not a pointer to struct"),
			},
		)
	})
	t.Run("string config object", func(t *testing.T) {
		var errLoad *LoadError
		require.ErrorAs(t,
			NewConfig(NewFakeSourcer("app", nil)).Load("foo"),
			&errLoad,
		)
		assert.ElementsMatch(t,
			errLoad.Errors,
			[]error{
				fmt.Errorf("configuration target is not a pointer to struct"),
			},
		)
	})
}

func TestConfigEmbeddedConfig(t *testing.T) {
	config := NewConfig(NewFakeSourcer("app", map[string]string{
		"APP_A": "1",
		"APP_B": "2",
		"APP_C": "3",
		"APP_X": "4",
		"APP_Y": "5",
	}))

	chunk := &testParentConfig{}
	require.Nil(t, config.Load(chunk))
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

	chunk := &testParentConfig{}
	require.Nil(t, config.Load(chunk, NewEnvTagPrefixer("foo")))
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

	chunk := &testParentConfig{}

	var postErr *PostLoadError
	require.ErrorAs(t, config.Load(chunk), &postErr)
	assert.EqualError(t,
		postErr.Unwrap(),
		"fields must be increasing",
	)
}

type testParentConfig struct {
	TestChildConfig
	X int `env:"x"`
	Y int `env:"y"`
}

type TestChildConfig struct {
	A int `env:"a"`
	B int `env:"b"`
	C int `env:"c"`
}

func (c *TestChildConfig) PostLoad() error {
	if c.A >= c.B || c.B >= c.C {
		return fmt.Errorf("fields must be increasing")
	}

	return nil
}

func TestConfigBadEmbeddedObjectType(t *testing.T) {
	type C struct {
		*TestChildConfig
		X int `env:"x"`
		Y int `env:"y"`
	}

	config := NewConfig(NewFakeSourcer("app", nil))
	chunk := &C{}

	var loadErr *LoadError
	require.ErrorAs(t, config.Load(chunk), &loadErr)
	assert.ElementsMatch(t,
		loadErr.Errors,
		[]error{
			fmt.Errorf("invalid embedded type in configuration struct"),
		},
	)
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
