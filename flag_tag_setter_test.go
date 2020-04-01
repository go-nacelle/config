package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FlagTagSetterSuite struct{}

func (s *FlagTagSetterSuite) TestFlagTagSetter(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&BasicConfig{},
		NewFlagTagSetter(),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"flag":    "a",
		"default": "q",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{}))
}

func (s *FlagTagSetterSuite) TestFlagTagSetterEmbedded(t sweet.T) {
	obj, err := ApplyTagModifiers(
		&ParentConfig{},
		NewFlagTagSetter(),
	)

	Expect(err).To(BeNil())

	Expect(gatherTags(obj, "X")).To(Equal(map[string]string{
		"env":     "a",
		"flag":    "a",
		"default": "q",
	}))

	Expect(gatherTags(obj, "Y")).To(Equal(map[string]string{}))
}
