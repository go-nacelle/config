package zubrin

import (
	"github.com/aphistic/sweet"
)

type DirectorySourcerSuite struct{}

func (s *DirectorySourcerSuite) TestLoadJSON(t sweet.T) {
	sourcer := NewDirectorySourcer("test-files/dir", nil)

	ensureEquals(sourcer, []string{"a"}, "1")
	ensureEquals(sourcer, []string{"b"}, "10")
	ensureEquals(sourcer, []string{"c"}, "100")
	ensureEquals(sourcer, []string{"d"}, "200")
	ensureEquals(sourcer, []string{"e"}, "300")
	ensureEquals(sourcer, []string{"x"}, "7")
	ensureEquals(sourcer, []string{"y"}, "8")
	ensureEquals(sourcer, []string{"z"}, "9")
}

func (s *DirectorySourcerSuite) TestOptionalDirectorySourcer(t sweet.T) {
	ensureMissing(
		NewOptionalFileSourcer("test-files/no-such-directory", nil),
		[]string{"foo"},
	)
}
