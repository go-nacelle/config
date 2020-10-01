package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagTagPrefixer(t *testing.T) {
	obj, err := ApplyTagModifiers(&BasicFlagConfig{}, NewFlagTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"flag": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}

func TestFlagTagPrefixerEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(&ParentFlagConfig{}, NewFlagTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"flag": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}
