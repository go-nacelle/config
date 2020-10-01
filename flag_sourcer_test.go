package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagSourcerGet(t *testing.T) {
	sourcer := NewFlagSourcer(WithFlagSourcerArgs([]string{"-X=foo", "--Y", "123"}))
	require.Nil(t, sourcer.Init())

	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"Y"})
	assert.Equal(t, "foo", val1)
	assert.Equal(t, "123", val2)
}

func TestFlagSourcerIllegalFlag(t *testing.T) {
	for _, badFlag := range []string{"X", "---X", "-=", "--="} {
		sourcer := NewFlagSourcer(WithFlagSourcerArgs([]string{badFlag}))
		assert.EqualError(t, sourcer.Init(), fmt.Sprintf("illegal flag: %s", badFlag))
	}
}

func TestFlagSourcerMissingArgument(t *testing.T) {
	sourcer := NewFlagSourcer(WithFlagSourcerArgs([]string{"--X"}))
	assert.EqualError(t, sourcer.Init(), fmt.Sprintf("flag needs an argument: -X"))
}

func TestFlagSourcerDump(t *testing.T) {
	sourcer := NewFlagSourcer(WithFlagSourcerArgs([]string{"-X=foo", "--Y", "123"}))
	require.Nil(t, sourcer.Init())

	dump := sourcer.Dump()
	assert.Equal(t, "foo", dump["X"])
	assert.Equal(t, "123", dump["Y"])
}
