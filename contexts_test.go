package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithAndFromContext(t *testing.T) {
	t.Run("from context without value", func(t *testing.T) {
		ctx := context.Background()
		assert.Nil(t, FromContext(ctx))
	})
	t.Run("from context that has value", func(t *testing.T) {
		cfg := NewConfig(NewFakeSourcer("app", map[string]string{
			"APP_FOO": "bar",
		}))

		ctx := context.Background()
		ctx = WithContext(ctx, cfg)

		ctxCfg := FromContext(ctx)
		require.NotNil(t, ctxCfg)

		type TestConfig struct {
			Foo string `env:"foo"`
		}

		loadCfg := &TestConfig{}
		err := ctxCfg.Load(loadCfg)
		assert.NoError(t, err)
		assert.Equal(t, "bar", loadCfg.Foo)
	})
}

func TestLoadFromContext(t *testing.T) {
	type TestConfig struct {
		Foo string `env:"foo"`
	}

	t.Run("no config in context", func(t *testing.T) {
		ctx := context.Background()

		loadCfg := &TestConfig{}
		err := LoadFromContext(ctx, loadCfg)
		assert.EqualError(t, err, "config loader not found in context")
	})
	t.Run("loads from context", func(t *testing.T) {
		cfg := NewConfig(NewFakeSourcer("app", map[string]string{
			"APP_FOO": "bar",
		}))

		ctx := context.Background()
		ctx = WithContext(ctx, cfg)

		loadCfg := &TestConfig{}
		err := LoadFromContext(ctx, loadCfg)
		assert.NoError(t, err)
		assert.Equal(t, "bar", loadCfg.Foo)
	})
	t.Run("loads from context with tag modifiers", func(t *testing.T) {
		cfg := NewConfig(NewFakeSourcer("app", map[string]string{
			"APP_TAG_FOO": "bar",
		}))

		ctx := context.Background()
		ctx = WithContext(ctx, cfg)

		loadCfg := &TestConfig{}
		err := LoadFromContext(ctx, loadCfg, NewEnvTagPrefixer("tag"))
		assert.NoError(t, err)
		assert.Equal(t, "bar", loadCfg.Foo)
	})
}
