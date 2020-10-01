package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestEnvSourcerUnprefixed(t *testing.T) {
	values := map[string]string{
		"X": "foo",
		"Y": "123",
	}

	sourcer := NewTestEnvSourcer(values)
	require.Nil(t, sourcer.Init())

	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"Y"})
	assert.Equal(t, "foo", val1)
	assert.Equal(t, "123", val2)
}

func TestTestEnvSourcerNormalizedCasing(t *testing.T) {
	values := map[string]string{
		"x": "foo",
		"y": "123",
	}

	sourcer := NewTestEnvSourcer(values)
	require.Nil(t, sourcer.Init())

	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"y"})
	assert.Equal(t, "foo", val1)
	assert.Equal(t, "123", val2)
}

func TestTestEnvSourcerDump(t *testing.T) {
	values := map[string]string{
		"X": "foo",
		"Y": "123",
	}

	sourcer := NewTestEnvSourcer(values)
	require.Nil(t, sourcer.Init())

	dump := sourcer.Dump()
	assert.Equal(t, "foo", dump["X"])
	assert.Equal(t, "123", dump["Y"])
}
