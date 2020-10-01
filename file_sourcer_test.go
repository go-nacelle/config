package config

import (
	"testing"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
)

func TestFileSourcerLoadJSON(t *testing.T) {
	sourcer := NewFileSourcer("test-files/values.json", ParseYAML)
	assert.Nil(t, sourcer.Init())
	testFileSourcer(t, sourcer)
}

func TestFileSourcerLoadJSONNoParser(t *testing.T) {
	sourcer := NewFileSourcer("test-files/values.json", nil)
	assert.Nil(t, sourcer.Init())
	testFileSourcer(t, sourcer)
}

func TestFileSourcerLoadYAML(t *testing.T) {
	sourcer := NewFileSourcer("test-files/values.yaml", ParseYAML)
	assert.Nil(t, sourcer.Init())
	testFileSourcer(t, sourcer)
}
func TestFileSourcerLoadYAMLNoParser(t *testing.T) {
	sourcer := NewFileSourcer("test-files/values.yaml", nil)
	assert.Nil(t, sourcer.Init())
	testFileSourcer(t, sourcer)
}

func TestFileSourcerLoadTOML(t *testing.T) {
	sourcer := NewFileSourcer("test-files/values.toml", ParseTOML)
	assert.Nil(t, sourcer.Init())
	testFileSourcer(t, sourcer)
}

func TestFileSourcerLoadTOMLNoParser(t *testing.T) {
	sourcer := NewFileSourcer("test-files/values.toml", nil)
	assert.Nil(t, sourcer.Init())
	testFileSourcer(t, sourcer)
}

func TestFileSourcerLoadJSONWithFakeFS(t *testing.T) {
	fs := NewMockFileSystem()
	fs.ReadFileFunc.SetDefaultReturn([]byte(`{
		"foo": "bar",
		"bar": [1, 2, 3],
		"baz": null,
		"bonk": {
			"x": 1,
			"y": 2,
			"z": 3
		},
		"encoded": "{\"w\": 4}",
		"deeply": {
			"nested": {
				"struct": [1, 2, 3]
			}
		}
	}`), nil)

	sourcer := NewFileSourcer("test-files/values.json", ParseYAML, WithFileSourcerFS(fs))
	assert.Nil(t, sourcer.Init())

	testFileSourcer(t, sourcer)
	mockassert.CalledOnceWith(t, fs.ReadFileFunc, mockassert.Values("test-files/values.json"))
}

func TestOptionalFileSourcer(t *testing.T) {
	sourcer := NewOptionalFileSourcer("test-files/no-such-file.json", nil)
	assert.Nil(t, sourcer.Init())
	ensureMissing(t, sourcer, []string{"foo"})
}

func TestOptionalFileSourcerWithFakeFS(t *testing.T) {
	fs := NewMockFileSystem()
	fs.ExistsFunc.SetDefaultReturn(false, nil)
	sourcer := NewOptionalFileSourcer("test-files/no-such-file.json", nil, WithFileSourcerFS(fs))
	assert.Nil(t, sourcer.Init())

	ensureMissing(t, sourcer, []string{"foo"})
	mockassert.CalledOnceWith(t, fs.ExistsFunc, mockassert.Values("test-files/no-such-file.json"))
}

func TestFileSourcerDump(t *testing.T) {
	sourcer := NewOptionalFileSourcer("test-files/values.json", ParseYAML)
	expected := map[string]string{
		"foo":     `bar`,
		"bar":     `[1,2,3]`,
		"baz":     `null`,
		"bonk":    `{"x":1,"y":2,"z":3}`,
		"encoded": `{"w": 4}`,
		"deeply":  `{"nested":{"struct":[1,2,3]}}`,
	}

	assert.Nil(t, sourcer.Init())
	assert.Equal(t, expected, sourcer.Dump())
}

func testFileSourcer(t *testing.T, sourcer Sourcer) {
	ensureEquals(t, sourcer, []string{"foo"}, "bar")
	ensureMatches(t, sourcer, []string{"bar"}, "[1, 2, 3]")
	ensureMatches(t, sourcer, []string{"bonk"}, `{"x": 1, "y": 2, "z": 3}`)
	ensureMatches(t, sourcer, []string{"encoded"}, `{"w": 4}`)
	ensureMatches(t, sourcer, []string{"bonk.x"}, `1`)
	ensureMatches(t, sourcer, []string{"encoded.w"}, `4`)
	ensureMatches(t, sourcer, []string{"deeply.nested.struct"}, `[1, 2, 3]`)
}
