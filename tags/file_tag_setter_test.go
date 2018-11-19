package tags

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/zubrin/internal/fixtures"
)

type FileTagSetterSuite struct{}

func (s *FileTagSetterSuite) TestFileTagSetter(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&fixtures.BasicConfig{},
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
		&fixtures.ParentConfig{},
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
