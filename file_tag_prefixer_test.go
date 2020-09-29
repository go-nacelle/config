package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTagPrefixer(t *testing.T) {
	obj, err := ApplyTagModifiers(&BasicFileConfig{}, NewFileTagPrefixer("foo"))
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"file": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}

func TestFileTagPrefixerEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(&ParentFileConfig{}, NewFileTagPrefixer("foo"))
	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"file": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}
