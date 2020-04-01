package config

import (
	"github.com/aphistic/sweet"
	. "github.com/efritz/go-mockgen/matchers"
	. "github.com/onsi/gomega"
)

type MultiSourcerSuite struct{}

func (s *MultiSourcerSuite) TestMultiSourcerBasic(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s1.TagsFunc.SetDefaultReturn([]string{"env"})
	s2.TagsFunc.SetDefaultReturn([]string{"env"})

	s1.GetFunc.SetDefaultHook(func(values []string) (string, SourcerFlag, error) {
		if values[0] == "foo" {
			return "bar", FlagFound, nil
		}

		return "", FlagMissing, nil
	})

	s2.GetFunc.SetDefaultHook(func(values []string) (string, SourcerFlag, error) {
		if values[0] == "bar" {
			return "baz", FlagFound, nil
		}

		return "", FlagMissing, nil
	})

	multi := NewMultiSourcer(s2, s1)
	ensureEquals(multi, []string{"foo"}, "bar")
	ensureEquals(multi, []string{"bar"}, "baz")
	ensureMissing(multi, []string{"baz"})
}

func (s *MultiSourcerSuite) TestMultiSourcerPriority(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s1.TagsFunc.SetDefaultReturn([]string{"env"})
	s2.TagsFunc.SetDefaultReturn([]string{"env"})
	s1.GetFunc.SetDefaultReturn("bar", FlagFound, nil)
	s2.GetFunc.SetDefaultReturn("baz", FlagFound, nil)

	multi := NewMultiSourcer(s2, s1)
	ensureEquals(multi, []string{"foo"}, "bar")
}

func (s *MultiSourcerSuite) TestMultiSourcerTags(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s4 := NewMockSourcer()
	s5 := NewMockSourcer()
	s1.TagsFunc.SetDefaultReturn([]string{"a"})
	s2.TagsFunc.SetDefaultReturn([]string{"b"})
	s3.TagsFunc.SetDefaultReturn([]string{"c"})
	s4.TagsFunc.SetDefaultReturn([]string{"a", "b", "d"})
	s5.TagsFunc.SetDefaultReturn([]string{"e"})

	multi := NewMultiSourcer(s5, s4, s3, s2, s1)
	tags := multi.Tags()
	Expect(tags).To(HaveLen(5))
	Expect(tags).To(ConsistOf("a", "b", "c", "d", "e"))
}

func (s *MultiSourcerSuite) TestMultiSourcerDifferentTags(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.TagsFunc.SetDefaultReturn([]string{"a"})
	s2.TagsFunc.SetDefaultReturn([]string{"b"})
	s3.TagsFunc.SetDefaultReturn([]string{"a"})
	s1.GetFunc.SetDefaultReturn("", FlagSkip, nil)
	s2.GetFunc.SetDefaultReturn("", FlagSkip, nil)
	s3.GetFunc.SetDefaultReturn("", FlagMissing, nil)

	multi := NewMultiSourcer(s3, s2, s1)
	_, flag, err := multi.Get([]string{"foo", "bar"})
	Expect(err).To(BeNil())
	Expect(flag).To(Equal(FlagMissing))
	Expect(s1.GetFunc).To(BeCalledOnceWith([]string{"foo"}))
	Expect(s2.GetFunc).To(BeCalledOnceWith([]string{"bar"}))
	Expect(s3.GetFunc).To(BeCalledOnceWith([]string{"foo"}))
}

func (s *MultiSourcerSuite) TestMultiSourceSkip(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.TagsFunc.SetDefaultReturn([]string{"a"})
	s2.TagsFunc.SetDefaultReturn([]string{"b"})
	s3.TagsFunc.SetDefaultReturn([]string{"a"})

	s1.GetFunc.SetDefaultReturn("", FlagSkip, nil)
	s2.GetFunc.SetDefaultReturn("", FlagSkip, nil)
	s3.GetFunc.SetDefaultReturn("", FlagSkip, nil)

	multi := NewMultiSourcer(s3, s2, s1)
	_, flag, err := multi.Get([]string{"", ""})
	Expect(err).To(BeNil())
	Expect(flag).To(Equal(FlagSkip))
}

func (s *MultiSourcerSuite) TestDump(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.DumpFunc.SetDefaultReturn(map[string]string{"a": "foo"}, nil)
	s2.DumpFunc.SetDefaultReturn(map[string]string{"b": "bar", "a": "bonk"}, nil)
	s3.DumpFunc.SetDefaultReturn(map[string]string{"c": "baz"}, nil)

	dump, err := NewMultiSourcer(s3, s2, s1).Dump()
	Expect(err).To(BeNil())
	Expect(dump).To(Equal(map[string]string{
		"a": "bonk",
		"b": "bar",
		"c": "baz",
	}))
}
