package config

import (
	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"
)

type DirectorySourcerSuite struct{}

func (s *DirectorySourcerSuite) TestLoadJSON(t sweet.T) {
	sourcer := NewDirectorySourcer("test-files/dir", nil)
	Expect(sourcer.Init()).To(BeNil())

	ensureEquals(sourcer, []string{"a"}, "1")
	ensureEquals(sourcer, []string{"b"}, "10")
	ensureEquals(sourcer, []string{"c"}, "100")
	ensureEquals(sourcer, []string{"d"}, "200")
	ensureEquals(sourcer, []string{"e"}, "300")
	ensureEquals(sourcer, []string{"x"}, "7")
	ensureEquals(sourcer, []string{"y"}, "8")
	ensureEquals(sourcer, []string{"z"}, "9")
}

func (s *DirectorySourcerSuite) TestLoadJSONWithFakeFS(t sweet.T) {
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
	Expect(sourcer.Init()).To(BeNil())

	ensureEquals(sourcer, []string{"a"}, "1")
	ensureEquals(sourcer, []string{"b"}, "10")
	ensureEquals(sourcer, []string{"c"}, "100")
	ensureEquals(sourcer, []string{"d"}, "200")
	ensureEquals(sourcer, []string{"e"}, "300")
	ensureEquals(sourcer, []string{"x"}, "7")
	ensureEquals(sourcer, []string{"y"}, "8")
	ensureEquals(sourcer, []string{"z"}, "9")

	Expect(fs.ListFilesFunc).To(BeCalledOnceWith("test-files/dir"))
	Expect(fs.ReadFileFunc).To(BeCalledOnceWith("test-files/dir/a.json"))
	Expect(fs.ReadFileFunc).To(BeCalledOnceWith("test-files/dir/b.json"))
	Expect(fs.ReadFileFunc).To(BeCalledOnceWith("test-files/dir/c.json"))
}

func (s *DirectorySourcerSuite) TestOptionalDirectorySourcer(t sweet.T) {
	sourcer := NewOptionalDirectorySourcer("test-files/no-such-directory", nil)
	Expect(sourcer.Init()).To(BeNil())
	ensureMissing(sourcer, []string{"foo"})
}

func (s *DirectorySourcerSuite) TestOptionalDirectorySourcerWithFakeFS(t sweet.T) {
	fs := NewMockFileSystem()
	fs.ExistsFunc.SetDefaultReturn(false, nil)

	sourcer := NewOptionalDirectorySourcer("test-files/no-such-directory", nil, WithDirectorySourcerFS(fs))
	Expect(sourcer.Init()).To(BeNil())

	ensureMissing(sourcer, []string{"foo"})
	Expect(fs.ExistsFunc).To(BeCalledOnceWith("test-files/no-such-directory"))
}
