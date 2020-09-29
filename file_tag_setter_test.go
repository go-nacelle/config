package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileTagSetter(t *testing.T) {
	obj, err := ApplyTagModifiers(
		&BasicConfig{},
		NewFileTagSetter(),
	)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "file": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}

func TestFileTagSetterEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(
		&ParentConfig{},
		NewFileTagSetter(),
	)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "file": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}
