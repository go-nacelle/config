package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTagSetter(t *testing.T) {
	obj, err := ApplyTagModifiers(
		&BasicConfig{},
		NewDefaultTagSetter("X", "r"),
		NewDefaultTagSetter("Y", "null"),
	)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "default": "r"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{"default": "null"}, gatherTags(obj, "Y"))
}

func TestDefaultTagSetterEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(
		&ParentConfig{},
		NewDefaultTagSetter("X", "r"),
		NewDefaultTagSetter("Y", "null"),
	)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "default": "r"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{"default": "null"}, gatherTags(obj, "Y"))
}
