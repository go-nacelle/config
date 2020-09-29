package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplayTagSetter(t *testing.T) {
	obj, err := ApplyTagModifiers(
		&BasicConfig{},
		NewDisplayTagSetter(),
	)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "display": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}

func TestDisplayTagSetterEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(
		&ParentConfig{},
		NewDisplayTagSetter(),
	)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "display": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}
