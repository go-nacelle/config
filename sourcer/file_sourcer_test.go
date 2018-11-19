package sourcer

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FileSourcerSuite struct{}

func (s *FileSourcerSuite) TestLoadJSON(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.json", parseYAML))
}

func (s *FileSourcerSuite) TestLoadJSONNoParser(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.json", nil))
}

func (s *FileSourcerSuite) TestLoadYAML(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.yaml", parseYAML))
}
func (s *FileSourcerSuite) TestLoadYAMLNoParser(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.yaml", nil))
}

func (s *FileSourcerSuite) TestLoadTOML(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.toml", parseTOML))
}

func (s *FileSourcerSuite) TestLoadTOMLNoParser(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.toml", nil))
}

func (s *FileSourcerSuite) TestOptionalFileSourcer(t sweet.T) {
	ensureMissing(
		NewOptionalFileSourcer("test-files/no-such-file.json", nil),
		[]string{"foo"},
	)
}
func (s *FileSourcerSuite) TestDump(t sweet.T) {
	sourcer := NewOptionalFileSourcer("test-files/values.json", parseYAML)

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
