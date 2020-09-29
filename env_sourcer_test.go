package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvSourcerUnprefixed(t *testing.T) {
	os.Setenv("X", "foo")
	os.Setenv("Y", "123")
	os.Setenv("APP_Y", "456")
	sourcer := NewEnvSourcer("app")
	assert.Nil(t, sourcer.Init())

	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"Y"})
	assert.Equal(t, "foo", val1)
	assert.Equal(t, "456", val2)
}

func TestEnvSourcerNormalizedPrefix(t *testing.T) {
	os.Setenv("FOO_BAR_X", "foo")
	os.Setenv("FOO_BAR_Y", "123")
	sourcer := NewEnvSourcer("$foo-^-bar@")
	assert.Nil(t, sourcer.Init())

	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"Y"})
	assert.Equal(t, "foo", val1)
	assert.Equal(t, "123", val2)
}

func TestEnvSourcerDump(t *testing.T) {
	os.Setenv("X", "foo")
	os.Setenv("Y", "123")
	sourcer := NewEnvSourcer("app")
	assert.Nil(t, sourcer.Init())

	dump := sourcer.Dump()
	assert.Equal(t, "foo", dump["X"])
	assert.Equal(t, "123", dump["Y"])
}
