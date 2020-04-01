package config

import (
	"fmt"

	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type FlagSourcerSuite struct{}

func (s *FlagSourcerSuite) TestGet(t sweet.T) {
	sourcer := NewFlagSourcer(WithFlagSourcerArgs([]string{"-X=foo", "--Y", "123"}))
	Expect(sourcer.Init()).To(BeNil())

	val1, _, _ := sourcer.Get([]string{"X"})
	val2, _, _ := sourcer.Get([]string{"Y"})
	Expect(val1).To(Equal("foo"))
	Expect(val2).To(Equal("123"))
}

func (s *FlagSourcerSuite) TestIllegalFlag(t sweet.T) {
	for _, badFlag := range []string{"X", "---X", "-=", "--="} {
		sourcer := NewFlagSourcer(WithFlagSourcerArgs([]string{badFlag}))
		Expect(sourcer.Init()).To(MatchError(fmt.Sprintf("illegal flag: %s", badFlag)))
	}
}

func (s *FlagSourcerSuite) TestMissingArgument(t sweet.T) {
	sourcer := NewFlagSourcer(WithFlagSourcerArgs([]string{"--X"}))
	Expect(sourcer.Init()).To(MatchError(fmt.Sprintf("flag needs an argument: -X")))
}

func (s *FlagSourcerSuite) TestDump(t sweet.T) {
	sourcer := NewFlagSourcer(WithFlagSourcerArgs([]string{"-X=foo", "--Y", "123"}))
	Expect(sourcer.Init()).To(BeNil())

	dump := sourcer.Dump()
	Expect(dump["X"]).To(Equal("foo"))
	Expect(dump["Y"]).To(Equal("123"))
}
