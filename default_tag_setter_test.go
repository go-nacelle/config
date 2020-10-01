package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultTagSetter(t *testing.T) {
	type C struct {
		X string `env:"a" default:"q"`
		Y string
	}

	obj, err := ApplyTagModifiers(&C{}, NewDefaultTagSetter("X", "r"), NewDefaultTagSetter("Y", "null"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "default": "r"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{"default": "null"}, gatherTags(obj, "Y"))
}

func TestDefaultTagSetterEmbedded(t *testing.T) {
	type C1 struct {
		X string `env:"a" default:"q"`
		Y string
	}
	type C2 struct{ C1 }

	obj, err := ApplyTagModifiers(&C2{}, NewDefaultTagSetter("X", "r"), NewDefaultTagSetter("Y", "null"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "default": "r"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{"default": "null"}, gatherTags(obj, "Y"))
}
