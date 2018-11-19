package tags

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"

	"github.com/efritz/zubrin/internal/fixtures"
)

type DisplayTagSetterSuite struct{}

func (s *DisplayTagSetterSuite) TestDisplayTagSetter(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&fixtures.BasicConfig{},
		NewDisplayTagSetter(),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"display": "a",
		"default": "q",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{}))
}

func (s *DisplayTagSetterSuite) TestDisplayTagSetterEmbedded(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&fixtures.ParentConfig{},
		NewDisplayTagSetter(),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"display": "a",
		"default": "q",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{}))
}
