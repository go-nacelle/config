package config

import (
	"testing"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/require"
)

func TestDirectorySourcerLoadJSON(t *testing.T) {
	sourcer := NewDirectorySourcer("test-files/dir", nil)
	require.Nil(t, sourcer.Init())

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
	require.Nil(t, sourcer.Init())

	ensureEquals(t, sourcer, []string{"a"}, "1")
	ensureEquals(t, sourcer, []string{"b"}, "10")
	ensureEquals(t, sourcer, []string{"c"}, "100")
	ensureEquals(t, sourcer, []string{"d"}, "200")
	ensureEquals(t, sourcer, []string{"e"}, "300")
	ensureEquals(t, sourcer, []string{"x"}, "7")
	ensureEquals(t, sourcer, []string{"y"}, "8")
	ensureEquals(t, sourcer, []string{"z"}, "9")

	mockassert.CalledOnceWith(t, fs.ListFilesFunc, mockassert.Values("test-files/dir"))
	mockassert.CalledOnceWith(t, fs.ReadFileFunc, mockassert.Values("test-files/dir/a.json"))
	mockassert.CalledOnceWith(t, fs.ReadFileFunc, mockassert.Values("test-files/dir/b.json"))
	mockassert.CalledOnceWith(t, fs.ReadFileFunc, mockassert.Values("test-files/dir/c.json"))
}

func TestOptionalDirectorySourcer(t *testing.T) {
	sourcer := NewOptionalDirectorySourcer("test-files/no-such-directory", nil)
	require.Nil(t, sourcer.Init())
	ensureMissing(t, sourcer, []string{"foo"})
}

func TestOptionalDirectorySourcerWithFakeFS(t *testing.T) {
	fs := NewMockFileSystem()
	fs.ExistsFunc.SetDefaultReturn(false, nil)

	sourcer := NewOptionalDirectorySourcer("test-files/no-such-directory", nil, WithDirectorySourcerFS(fs))
	require.Nil(t, sourcer.Init())

	ensureMissing(t, sourcer, []string{"foo"})
	mockassert.CalledOnceWith(t, fs.ExistsFunc, mockassert.Values("test-files/no-such-directory"))
}
