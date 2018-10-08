package zubrin

import "github.com/aphistic/sweet"

type FileSourcerSuite struct{}

func (s *FileSourcerSuite) TestLoadJSON(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.json", ParseYAML))
}

func (s *FileSourcerSuite) TestLoadJSONNoParser(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.json", nil))
}

func (s *FileSourcerSuite) TestLoadYAML(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.yaml", ParseYAML))
}
func (s *FileSourcerSuite) TestLoadYAMLNoParser(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.yaml", nil))
}

func (s *FileSourcerSuite) TestLoadTOML(t sweet.T) {
	testFileSourcer(NewFileSourcer("test-files/values.toml", ParseTOML))
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

func testFileSourcer(sourcer Sourcer) {
	ensureEquals(sourcer, []string{"foo"}, "bar")
	ensureMatches(sourcer, []string{"bar"}, "[1, 2, 3]")
	ensureMatches(sourcer, []string{"bonk"}, `{"x": 1, "y": 2, "z": 3}`)
	ensureMatches(sourcer, []string{"encoded"}, `{"w": 4}`)
	ensureMatches(sourcer, []string{"bonk.x"}, `1`)
	ensureMatches(sourcer, []string{"encoded.w"}, `4`)
	ensureMatches(sourcer, []string{"deeply.nested.struct"}, `[1, 2, 3]`)
}
