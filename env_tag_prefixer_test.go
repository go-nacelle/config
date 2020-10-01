package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvTagPrefixer(t *testing.T) {
	obj, err := ApplyTagModifiers(&BasicConfig{}, NewEnvTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}

func TestEnvTagPrefixerEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(&ParentConfig{}, NewEnvTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}
