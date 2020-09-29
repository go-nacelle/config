package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagTagPrefixer(t *testing.T) {
	obj, err := ApplyTagModifiers(&BasicFlagConfig{}, NewFlagTagPrefixer("foo"))
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"flag": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}

func TestFlagTagPrefixerEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(&ParentFlagConfig{}, NewFlagTagPrefixer("foo"))
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"flag": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}
