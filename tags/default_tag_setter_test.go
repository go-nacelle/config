package tags

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/zubrin/internal/fixtures"
)

type DefaultTagSetterSuite struct{}

func (s *DefaultTagSetterSuite) TestDefaultTagSetter(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&fixtures.BasicConfig{},
		NewDefaultTagSetter("X", "r"),
		NewDefaultTagSetter("Y", "null"),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"default": "r",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{
		"default": "null",
	}))
}

func (s *DefaultTagSetterSuite) TestDefaultTagSetterEmbedded(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&fixtures.ParentConfig{},
		NewDefaultTagSetter("X", "r"),
		NewDefaultTagSetter("Y", "null"),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"default": "r",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{
		"default": "null",
	}))
}
