package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagTagPrefixer(t *testing.T) {
	type C1 struct {
		X string `flag:"a" default:"q"`
		Y string
	}

	obj, err := ApplyTagModifiers(&C1{}, NewFlagTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"flag": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}

func TestFlagTagPrefixerEmbedded(t *testing.T) {
	type C1 struct {
		X string `flag:"a" default:"q"`
		Y string
	}
	type C2 struct{ C1 }

	obj, err := ApplyTagModifiers(&C2{}, NewFlagTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"flag": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}
