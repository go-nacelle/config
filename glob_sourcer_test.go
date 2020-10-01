package config

import (
	"testing"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
)

func TestGlobSourcerLoadJSON(t *testing.T) {
	sourcer := NewGlobSourcer("test-files/**/*.json", nil)
	assert.Nil(t, sourcer.Init())

	ensureEquals(t, sourcer, []string{"nested-x"}, "1")
	ensureEquals(t, sourcer, []string{"nested-y"}, "2")
	ensureEquals(t, sourcer, []string{"nested-z"}, "3")
	ensureEquals(t, sourcer, []string{"nested-w"}, "4")
}

func TestGlobSourcerLoadJSONWithFakeFS(t *testing.T) {
	fs := NewMockFileSystem()
	fs.GlobFunc.SetDefaultReturn([]string{
		"test-files/dir/nested-a/x.json",
		"test-files/dir/nested-b/y.json",
		"test-files/dir/nested-b/z.json",
		"test-files/dir/nested-b/nested-c/w.json",
	}, nil)

	fs.ReadFileFunc.PushReturn([]byte(`{"nested-x": 1}`), nil)
	fs.ReadFileFunc.PushReturn([]byte(`{"nested-y": 2}`), nil)
	fs.ReadFileFunc.PushReturn([]byte(`{"nested-z": 3}`), nil)
	fs.ReadFileFunc.PushReturn([]byte(`{"nested-w": 4}`), nil)

	sourcer := NewGlobSourcer("test-files/**/*.json", nil, WithGlobSourcerFS(fs))
	assert.Nil(t, sourcer.Init())

	mockassert.CalledOnceWith(t, fs.GlobFunc, mockassert.Values("test-files/**/*.json"))
	mockassert.CalledOnceWith(t, fs.ReadFileFunc, mockassert.Values("test-files/dir/nested-a/x.json"))
	mockassert.CalledOnceWith(t, fs.ReadFileFunc, mockassert.Values("test-files/dir/nested-b/y.json"))
	mockassert.CalledOnceWith(t, fs.ReadFileFunc, mockassert.Values("test-files/dir/nested-b/z.json"))
	mockassert.CalledOnceWith(t, fs.ReadFileFunc, mockassert.Values("test-files/dir/nested-b/nested-c/w.json"))

	ensureEquals(t, sourcer, []string{"nested-x"}, "1")
	ensureEquals(t, sourcer, []string{"nested-y"}, "2")
	ensureEquals(t, sourcer, []string{"nested-z"}, "3")
	ensureEquals(t, sourcer, []string{"nested-w"}, "4")
}

func TestGlobSourcerNoMatches(t *testing.T) {
	sourcer := NewGlobSourcer("test-files/notexist/*.yaml", nil)
	assert.Nil(t, sourcer.Init())
	assert.Empty(t, sourcer.Tags())
}
