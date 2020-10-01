package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileTagPrefixer(t *testing.T) {
	type C struct {
		X string `file:"a" default:"q"`
		Y string
	}

	obj, err := ApplyTagModifiers(&C{}, NewFileTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"file": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}

func TestFileTagPrefixerEmbedded(t *testing.T) {
	type C1 struct {
		X string `file:"a" default:"q"`
		Y string
	}
	type C2 struct{ C1 }

	obj, err := ApplyTagModifiers(&C2{}, NewFileTagPrefixer("foo"))
	require.Nil(t, err)
	assert.Equal(t, map[string]string{"file": "foo_a", "default": "q"}, gatherTags(obj, "X"))
}
