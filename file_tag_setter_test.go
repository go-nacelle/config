package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileTagSetter(t *testing.T) {
	type C struct {
		X string `env:"a" default:"q"`
		Y string
	}

	obj, err := ApplyTagModifiers(&C{}, NewFileTagSetter())
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "file": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}

func TestFileTagSetterEmbedded(t *testing.T) {
	type C1 struct {
		X string `env:"a" default:"q"`
		Y string
	}
	type C2 struct{ C1 }

	obj, err := ApplyTagModifiers(&C2{}, NewFileTagSetter())
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "file": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}
