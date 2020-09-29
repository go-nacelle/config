package config

import (
	"testing"

	mockassert "github.com/efritz/go-mockgen/assert"
	"github.com/stretchr/testify/assert"
)

func TestDirectorySourcerLoadJSON(t *testing.T) {
	sourcer := NewDirectorySourcer("test-files/dir", nil)
	assert.Nil(t, sourcer.Init())

	ensureEquals(t, sourcer, []string{"a"}, "1")
	ensureEquals(t, sourcer, []string{"b"}, "10")
	ensureEquals(t, sourcer, []string{"c"}, "100")
	ensureEquals(t, sourcer, []string{"d"}, "200")
	ensureEquals(t, sourcer, []string{"e"}, "300")
	ensureEquals(t, sourcer, []string{"x"}, "7")
	ensureEquals(t, sourcer, []string{"y"}, "8")
	ensureEquals(t, sourcer, []string{"z"}, "9")
}

func TestDirectorySourcerLoadJSONWithFakeFS(t *testing.T) {
	fs := NewMockFileSystem()
	fs.ListFilesFunc.SetDefaultReturn([]string{"a.json", "b.json", "c.json"}, nil)

	fs.ReadFileFunc.PushReturn([]byte(`{
		"a": 1,
		"b": 2,
		"c": 3,
		"x": 7
	}`), nil)

	fs.ReadFileFunc.PushReturn([]byte(`{
		"b": 10,
		"c": 20,
		"d": 30,
		"y": 8
	}`), nil)

	fs.ReadFileFunc.PushReturn([]byte(`{
		"c": 100,
		"d": 200,
		"e": 300,
		"z": 9
	}`), nil)

	sourcer := NewDirectorySourcer("test-files/dir", nil, WithDirectorySourcerFS((fs)))
	assert.Nil(t, sourcer.Init())

	ensureEquals(t, sourcer, []string{"a"}, "1")
	ensureEquals(t, sourcer, []string{"b"}, "10")
	ensureEquals(t, sourcer, []string{"c"}, "100")
	ensureEquals(t, sourcer, []string{"d"}, "200")
	ensureEquals(t, sourcer, []string{"e"}, "300")
	ensureEquals(t, sourcer, []string{"x"}, "7")
	ensureEquals(t, sourcer, []string{"y"}, "8")
	ensureEquals(t, sourcer, []string{"z"}, "9")

	mockassert.CalledOnceMatching(t, fs.ListFilesFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemListFilesFuncCall).Arg0 == "test-files/dir" // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, fs.ReadFileFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemReadFileFuncCall).Arg0 == "test-files/dir/a.json" // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, fs.ReadFileFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemReadFileFuncCall).Arg0 == "test-files/dir/b.json" // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, fs.ReadFileFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemReadFileFuncCall).Arg0 == "test-files/dir/c.json" // TODO - ergonomics
	})
}

func TestOptionalDirectorySourcer(t *testing.T) {
	sourcer := NewOptionalDirectorySourcer("test-files/no-such-directory", nil)
	assert.Nil(t, sourcer.Init())
	ensureMissing(t, sourcer, []string{"foo"})
}

func TestOptionalDirectorySourcerWithFakeFS(t *testing.T) {
	fs := NewMockFileSystem()
	fs.ExistsFunc.SetDefaultReturn(false, nil)

	sourcer := NewOptionalDirectorySourcer("test-files/no-such-directory", nil, WithDirectorySourcerFS(fs))
	assert.Nil(t, sourcer.Init())

	ensureMissing(t, sourcer, []string{"foo"})
	mockassert.CalledOnceMatching(t, fs.ExistsFunc, func(t assert.TestingT, call interface{}) bool {
		return call.(FileSystemExistsFuncCall).Arg0 == "test-files/no-such-directory" // TODO - ergonomics
	})
}
