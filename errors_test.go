package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadError(t *testing.T) {
	t.Run("error string includes error slice", func(t *testing.T) {
		err := newLoadError([]error{
			fmt.Errorf("err1"),
			fmt.Errorf("err2"),
		})

		assert.EqualError(t,
			err,
			"failed to load config: err1, err2",
		)
	})
}

func TestSerializeError(t *testing.T) {
	t.Run("unwrap error", func(t *testing.T) {
		innerErr := fmt.Errorf("inner")
		err := newSerializeError("X", innerErr)
		assert.Equal(t, "X", err.FieldName)
		assert.ErrorIs(t, err, innerErr)
	})
	t.Run("error string", func(t *testing.T) {
		innerErr := fmt.Errorf("inner")
		err := newSerializeError("X", innerErr)
		assert.EqualError(t,
			err,
			"field 'X': inner",
		)
	})
}

func TestPostLoadError(t *testing.T) {
	t.Run("error string", func(t *testing.T) {
		err := newPostLoadError(fmt.Errorf("inner err"))
		assert.EqualError(t,
			err,
			"post load callback failed: inner err",
		)
	})
	t.Run("unwrap error", func(t *testing.T) {
		err := newPostLoadError(fmt.Errorf("inner err"))

		innerErr := err.Unwrap()
		require.NotNil(t, innerErr)
		assert.EqualError(t, innerErr, "inner err")
	})
}
