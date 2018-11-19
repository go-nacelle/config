package zubrin

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FileTagSetterSuite struct{}

func (s *FileTagSetterSuite) TestFileTagSetter(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&BasicConfig{},
		NewFileTagSetter(),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"file":    "a",
		"default": "q",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{}))
}

func (s *FileTagSetterSuite) TestFileTagSetterEmbedded(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&ParentConfig{},
		NewFileTagSetter(),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"file":    "a",
		"default": "q",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{}))
}
