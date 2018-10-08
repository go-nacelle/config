package zubrin

//go:generate go-mockgen github.com/efritz/zubrin -i Sourcer -o mock_sourcer_test.go -f

import (
	"github.com/aphistic/sweet"
	. "github.com/onsi/gomega"
)

type MultiSourcerSuite struct{}

func (s *MultiSourcerSuite) TestMultiSourcerBasic(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"env"} }
	s2.TagsFunc = func() []string { return []string{"env"} }

	s1.GetFunc = func(values []string) (string, SourcerFlag, error) {
		if values[0] == "foo" {
			return "bar", FlagFound, nil
		}

		return "", FlagMissing, nil
	}

	s2.GetFunc = func(values []string) (string, SourcerFlag, error) {
		if values[0] == "bar" {
			return "baz", FlagFound, nil
		}

		return "", FlagMissing, nil
	}

	multi := NewMultiSourcer(s2, s1)
	ensureEquals(multi, []string{"foo"}, "bar")
	ensureEquals(multi, []string{"bar"}, "baz")
	ensureMissing(multi, []string{"baz"})
}

func (s *MultiSourcerSuite) TestMultiSourcerPriority(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"env"} }
	s2.TagsFunc = func() []string { return []string{"env"} }

	s1.GetFunc = func(values []string) (string, SourcerFlag, error) {
		return "bar", FlagFound, nil
	}

	s2.GetFunc = func(values []string) (string, SourcerFlag, error) {
		return "baz", FlagFound, nil
	}

	multi := NewMultiSourcer(s2, s1)
	ensureEquals(multi, []string{"foo"}, "bar")
}

func (s *MultiSourcerSuite) TestMultiSourcerTags(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s4 := NewMockSourcer()
	s5 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"a"} }
	s2.TagsFunc = func() []string { return []string{"b"} }
	s3.TagsFunc = func() []string { return []string{"c"} }
	s4.TagsFunc = func() []string { return []string{"a", "b", "d"} }
	s5.TagsFunc = func() []string { return []string{"e"} }

	multi := NewMultiSourcer(s5, s4, s3, s2, s1)
	tags := multi.Tags()
	Expect(tags).To(HaveLen(5))
	Expect(tags).To(ConsistOf("a", "b", "c", "d", "e"))
}

func (s *MultiSourcerSuite) TestMultiSourcerDifferentTags(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"a"} }
	s2.TagsFunc = func() []string { return []string{"b"} }
	s3.TagsFunc = func() []string { return []string{"a"} }

	s1.GetFunc = func(values []string) (string, SourcerFlag, error) {
		Expect(values).To(Equal([]string{"foo"}))
		return "", FlagSkip, nil
	}

	s2.GetFunc = func(values []string) (string, SourcerFlag, error) {
		Expect(values).To(Equal([]string{"bar"}))
		return "", FlagSkip, nil
	}

	s3.GetFunc = func(values []string) (string, SourcerFlag, error) {
		Expect(values).To(Equal([]string{"foo"}))
		return "", FlagMissing, nil
	}

	multi := NewMultiSourcer(s3, s2, s1)
	_, flag, err := multi.Get([]string{"foo", "bar"})
	Expect(err).To(BeNil())
	Expect(flag).To(Equal(FlagMissing))
	Expect(s1.GetFuncCallCount()).To(Equal(1))
	Expect(s2.GetFuncCallCount()).To(Equal(1))
	Expect(s3.GetFuncCallCount()).To(Equal(1))
}

func (s *MultiSourcerSuite) TestMultiSourceSkip(t sweet.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.TagsFunc = func() []string { return []string{"a"} }
	s2.TagsFunc = func() []string { return []string{"b"} }
	s3.TagsFunc = func() []string { return []string{"a"} }

	s1.GetFunc = func(values []string) (string, SourcerFlag, error) {
		return "", FlagSkip, nil
	}

	s2.GetFunc = func(values []string) (string, SourcerFlag, error) {
		return "", FlagSkip, nil
	}

	s3.GetFunc = func(values []string) (string, SourcerFlag, error) {
		return "", FlagSkip, nil
	}

	multi := NewMultiSourcer(s3, s2, s1)
	_, flag, err := multi.Get([]string{"", ""})
	Expect(err).To(BeNil())
	Expect(flag).To(Equal(FlagSkip))
}
