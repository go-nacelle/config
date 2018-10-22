package zubrin

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FileTagPrefixerSuite struct{}

func (s *FileTagPrefixerSuite) TestEnvTagPrefixer(t sweet.T) {
	obj, err := ApplyTagModifiers(&BasicConfig{}, NewEnvTagPrefixer("foo"))
	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "foo_a",
		"default": "q",
	}))
}

func (s *FileTagPrefixerSuite) TestEnvTagPrefixerEmbedded(t sweet.T) {
	obj, err := ApplyTagModifiers(&ParentConfig{}, NewEnvTagPrefixer("foo"))
	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "foo_a",
		"default": "q",
	}))
}
