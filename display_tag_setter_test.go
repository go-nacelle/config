package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisplayTagSetter(t *testing.T) {
	type C struct {
		X string `env:"a" default:"q"`
		Y string
	}

	obj, err := ApplyTagModifiers(&C{}, NewDisplayTagSetter())
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "display": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}

func TestDisplayTagSetterEmbedded(t *testing.T) {
	type C1 struct {
		X string `env:"a" default:"q"`
		Y string
	}
	type C2 struct{ C1 }

	obj, err := ApplyTagModifiers(&C2{}, NewDisplayTagSetter())
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"env": "a", "display": "a", "default": "q"}, gatherTags(obj, "X"))
	assert.Equal(t, map[string]string{}, gatherTags(obj, "Y"))
}
