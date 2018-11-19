package sourcer

import (
	"os"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type EnvSourcerSuite struct{}

func (s *EnvSourcerSuite) TestUnprefixed(t sweet.T) {
	os.Setenv("X", "foo")
	os.Setenv("Y", "123")
	os.Setenv("APP_Y", "456")

	sourcer := NewEnvSourcer("app")
	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"Y"})
	Expect(val1).To(Equal("foo"))
	Expect(val2).To(Equal("456"))
}

func (s *EnvSourcerSuite) TestNormalizedPrefix(t sweet.T) {
	os.Setenv("FOO_BAR_X", "foo")
	os.Setenv("FOO_BAR_Y", "123")

	sourcer := NewEnvSourcer("$foo-^-bar@")
	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"Y"})
	Expect(val1).To(Equal("foo"))
	Expect(val2).To(Equal("123"))
}

func (s *EnvSourcerSuite) TestDump(t sweet.T) {
	os.Setenv("X", "foo")
	os.Setenv("Y", "123")

	Expect(NewEnvSourcer("app").Dump()["X"]).To(Equal("foo"))
	Expect(NewEnvSourcer("app").Dump()["Y"]).To(Equal("123"))
}
