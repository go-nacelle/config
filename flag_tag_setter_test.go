package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlagTagSetter(t *testing.T) {
	obj, err := ApplyTagModifiers(
		&BasicConfig{},
		NewFlagTagSetter(),
	)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "flag": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}

func TestFlagTagSetterEmbedded(t *testing.T) {
	obj, err := ApplyTagModifiers(
		&ParentConfig{},
		NewFlagTagSetter(),
	)

	assert.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "flag": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}
