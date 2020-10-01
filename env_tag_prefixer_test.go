package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvTagPrefixer(t *testing.T) {
	type C struct {
		X string `env:"a" default:"q"`
		Y string
	}

	obj, err := ApplyTagModifiers(&C{}, NewEnvTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}

func TestEnvTagPrefixerEmbedded(t *testing.T) {
	type C1 struct {
		X string `env:"a" default:"q"`
		Y string
	}
	type C2 struct{ C1 }

	obj, err := ApplyTagModifiers(&C2{}, NewEnvTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}
