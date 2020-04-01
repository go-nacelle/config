package config

import (
	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"
)

type FileSourcerSuite struct{}

func (s *FileSourcerSuite) TestLoadJSON(t sweet.T) {
	sourcer := NewFileSourcer("test-files/values.json", ParseYAML)
	Expect(sourcer.Init()).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadJSONNoParser(t sweet.T) {
	sourcer := NewFileSourcer("test-files/values.json", nil)
	Expect(sourcer.Init()).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadYAML(t sweet.T) {
	sourcer := NewFileSourcer("test-files/values.yaml", ParseYAML)
	Expect(sourcer.Init()).To(BeNil())
	testFileSourcer(sourcer)
}
func (s *FileSourcerSuite) TestLoadYAMLNoParser(t sweet.T) {
	sourcer := NewFileSourcer("test-files/values.yaml", nil)
	Expect(sourcer.Init()).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadTOML(t sweet.T) {
	sourcer := NewFileSourcer("test-files/values.toml", ParseTOML)
	Expect(sourcer.Init()).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadTOMLNoParser(t sweet.T) {
	sourcer := NewFileSourcer("test-files/values.toml", nil)
	Expect(sourcer.Init()).To(BeNil())
	testFileSourcer(sourcer)
}

func (s *FileSourcerSuite) TestLoadJSONWithFakeFS(t sweet.T) {
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
	Expect(sourcer.Init()).To(BeNil())

	testFileSourcer(sourcer)
	Expect(fs.ReadFileFunc).To(BeCalledOnceWith("test-files/values.json"))
}

func (s *FileSourcerSuite) TestOptionalFileSourcer(t sweet.T) {
	sourcer := NewOptionalFileSourcer("test-files/no-such-file.json", nil)
	Expect(sourcer.Init()).To(BeNil())
	ensureMissing(sourcer, []string{"foo"})
}

func (s *FileSourcerSuite) TestOptionalFileSourcerWithFakeFS(t sweet.T) {
	fs := NewMockFileSystem()
	fs.ExistsFunc.SetDefaultReturn(false, nil)
	sourcer := NewOptionalFileSourcer("test-files/no-such-file.json", nil, WithFileSourcerFS(fs))
	Expect(sourcer.Init()).To(BeNil())

	ensureMissing(sourcer, []string{"foo"})
	Expect(fs.ExistsFunc).To(BeCalledOnceWith("test-files/no-such-file.json"))
}

func (s *FileSourcerSuite) TestDump(t sweet.T) {
	sourcer := NewOptionalFileSourcer("test-files/values.json", ParseYAML)
	Expect(sourcer.Init()).To(BeNil())

	Expect(sourcer.Dump()).To(Equal(map[string]string{
		"foo":     `bar`,
		"bar":     `[1,2,3]`,
		"baz":     `null`,
		"bonk":    `{"x":1,"y":2,"z":3}`,
		"encoded": `{"w": 4}`,
		"deeply":  `{"nested":{"struct":[1,2,3]}}`,
	}))
}

func testFileSourcer(sourcer Sourcer) {
	ensureEquals(sourcer, []string{"foo"}, "bar")
	ensureMatches(sourcer, []string{"bar"}, "[1, 2, 3]")
	ensureMatches(sourcer, []string{"bonk"}, `{"x": 1, "y": 2, "z": 3}`)
	ensureMatches(sourcer, []string{"encoded"}, `{"w": 4}`)
	ensureMatches(sourcer, []string{"bonk.x"}, `1`)
	ensureMatches(sourcer, []string{"encoded.w"}, `4`)
	ensureMatches(sourcer, []string{"deeply.nested.struct"}, `[1, 2, 3]`)
}
