package config

import (
	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"
)

type GlobSourcerSuite struct{}

func (s *GlobSourcerSuite) TestLoadJSON(t sweet.T) {
	sourcer := NewGlobSourcer("test-files/**/*.json", nil)

	ensureEquals(sourcer, []string{"nested-x"}, "1")
	ensureEquals(sourcer, []string{"nested-y"}, "2")
	ensureEquals(sourcer, []string{"nested-z"}, "3")
	ensureEquals(sourcer, []string{"nested-w"}, "4")
}

func (s *GlobSourcerSuite) TestLoadJSONWithFakeFS(t sweet.T) {
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
	Expect(fs.GlobFunc).To(BeCalledOnceWith("test-files/**/*.json"))
	Expect(fs.ReadFileFunc).To(BeCalledOnceWith("test-files/dir/nested-a/x.json"))
	Expect(fs.ReadFileFunc).To(BeCalledOnceWith("test-files/dir/nested-b/y.json"))
	Expect(fs.ReadFileFunc).To(BeCalledOnceWith("test-files/dir/nested-b/z.json"))
	Expect(fs.ReadFileFunc).To(BeCalledOnceWith("test-files/dir/nested-b/nested-c/w.json"))

	ensureEquals(sourcer, []string{"nested-x"}, "1")
	ensureEquals(sourcer, []string{"nested-y"}, "2")
	ensureEquals(sourcer, []string{"nested-z"}, "3")
	ensureEquals(sourcer, []string{"nested-w"}, "4")
}

func (s *GlobSourcerSuite) TestNoMatches(t sweet.T) {
	sourcer := NewGlobSourcer("test-files/notexist/*.yaml", nil)
	Expect(sourcer.Tags()).To(BeEmpty())
}
