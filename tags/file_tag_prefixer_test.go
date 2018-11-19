package tags

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/zubrin/internal/fixtures"
)

type FileTagPrefixerSuite struct{}

func (s *FileTagPrefixerSuite) TestEnvTagPrefixer(t sweet.T) {
	obj, err := ApplyTagModifiers(&fixtures.BasicConfig{}, NewEnvTagPrefixer("foo"))
	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "foo_a",
		"default": "q",
	}))
}

func (s *FileTagPrefixerSuite) TestEnvTagPrefixerEmbedded(t sweet.T) {
	obj, err := ApplyTagModifiers(&fixtures.ParentConfig{}, NewEnvTagPrefixer("foo"))
	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "foo_a",
		"default": "q",
	}))
}
