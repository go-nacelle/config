package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvTagPrefixer(t *testing.T) {
	obj, err := ApplyTagModifiers(&BasicConfig{}, NewEnvTagPrefixer("foo"))
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}

func TestEnvTagPrefixerEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(&ParentConfig{}, NewEnvTagPrefixer("foo"))
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}
