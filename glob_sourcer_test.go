package config

import (
	"testing"

	mockassert "github.com/efritz/go-mockgen/assert"
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

	mockassert.CalledOnceMatching(t, fs.GlobFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemGlobFuncCall).Arg0 == "test-files/**/*.json" // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, fs.ReadFileFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemReadFileFuncCall).Arg0 == "test-files/dir/nested-a/x.json" // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, fs.ReadFileFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemReadFileFuncCall).Arg0 == "test-files/dir/nested-b/y.json" // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, fs.ReadFileFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemReadFileFuncCall).Arg0 == "test-files/dir/nested-b/z.json" // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, fs.ReadFileFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemReadFileFuncCall).Arg0 == "test-files/dir/nested-b/nested-c/w.json" // TODO - ergonomics
	})

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
