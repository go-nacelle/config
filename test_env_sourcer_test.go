package config

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type TestEnvSourcerSuite struct{}

func (s *TestEnvSourcerSuite) TestUnprefixed(t sweet.T) {
	values := map[string]string{
		"X": "foo",
		"Y": "123",
	}

	sourcer := NewTestEnvSourcer(values)
	Expect(sourcer.Init()).To(BeNil())

	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"Y"})
	Expect(val1).To(Equal("foo"))
	Expect(val2).To(Equal("123"))
}

func (s *TestEnvSourcerSuite) TestNormalizedCasing(t sweet.T) {
	values := map[string]string{
		"x": "foo",
		"y": "123",
	}

	sourcer := NewTestEnvSourcer(values)
	Expect(sourcer.Init()).To(BeNil())

	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"y"})
	Expect(val1).To(Equal("foo"))
	Expect(val2).To(Equal("123"))
}

func (s *TestEnvSourcerSuite) TestDump(t sweet.T) {
	values := map[string]string{
		"X": "foo",
		"Y": "123",
	}

	sourcer := NewTestEnvSourcer(values)
	Expect(sourcer.Init()).To(BeNil())

	dump := sourcer.Dump()
	Expect(dump["X"]).To(Equal("foo"))
	Expect(dump["Y"]).To(Equal("123"))
}
