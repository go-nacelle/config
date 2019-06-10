package config

import (
	"github.com/aphistic/sweet"
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

func (s *GlobSourcerSuite) TestNoMatches(t sweet.T) {
	sourcer := NewGlobSourcer("test-files/notexist/*.yaml", nil)
	Expect(sourcer.Tags()).To(BeEmpty())
}
