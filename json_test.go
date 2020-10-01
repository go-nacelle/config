package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToJSONString(t *testing.T) {
	var val string
	require.True(t, toJSON([]byte("foobar"), &val))
	assert.Equal(t, "foobar", val)
}

func TestToJSONNonString(t *testing.T) {
	var val []int
	require.True(t, toJSON([]byte("[1, 2, 3, 4, 5]"), &val))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, val)
}

func TestToJSONBadType(t *testing.T) {
	var val []int
	assert.False(t, toJSON([]byte(`[1, 2, "3", 4, 5]`), &val))
}

func TestQuoteJSON(t *testing.T) {
	json := quoteJSON([]byte(`
	foo
	bar
	baz`))

	assert.JSONEq(t, `"\n\tfoo\n\tbar\n\tbaz"`, string(json))
}
