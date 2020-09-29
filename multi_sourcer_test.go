package config

import (
	"sort"
	"testing"

	mockassert "github.com/efritz/go-mockgen/assert"
	"github.com/stretchr/testify/assert"
)

func TestMultiSourcerBasic(t *testing.T) {
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

	sourcer := NewMultiSourcer(s2, s1)
	assert.Nil(t, sourcer.Init())

	ensureEquals(t, sourcer, []string{"foo"}, "bar")
	ensureEquals(t, sourcer, []string{"bar"}, "baz")
	ensureMissing(t, sourcer, []string{"baz"})
}

func TestMultiSourcerPriority(t *testing.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s1.TagsFunc.SetDefaultReturn([]string{"env"})
	s2.TagsFunc.SetDefaultReturn([]string{"env"})
	s1.GetFunc.SetDefaultReturn("bar", FlagFound, nil)
	s2.GetFunc.SetDefaultReturn("baz", FlagFound, nil)

	sourcer := NewMultiSourcer(s2, s1)
	assert.Nil(t, sourcer.Init())
	ensureEquals(t, sourcer, []string{"foo"}, "bar")
}

func TestMultiSourcerTags(t *testing.T) {
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

	sourcer := NewMultiSourcer(s5, s4, s3, s2, s1)
	assert.Nil(t, sourcer.Init())

	tags := sourcer.Tags()
	sort.Strings(tags)
	assert.Equal(t, []string{"a", "b", "c", "d", "e"}, tags)
}

func TestMultiSourcerDifferentTags(t *testing.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.TagsFunc.SetDefaultReturn([]string{"a"})
	s2.TagsFunc.SetDefaultReturn([]string{"b"})
	s3.TagsFunc.SetDefaultReturn([]string{"a"})
	s1.GetFunc.SetDefaultReturn("", FlagSkip, nil)
	s2.GetFunc.SetDefaultReturn("", FlagSkip, nil)
	s3.GetFunc.SetDefaultReturn("", FlagMissing, nil)

	sourcer := NewMultiSourcer(s3, s2, s1)
	assert.Nil(t, sourcer.Init())

	_, flag, err := sourcer.Get([]string{"foo", "bar"})
	assert.Nil(t, err)
	assert.Equal(t, FlagMissing, flag)

	mockassert.CalledOnceMatching(t, s1.GetFunc, func(t assert.TestingT, call interface{}) bool {
		return assert.Equal(t, []string{"foo"}, call.(SourcerGetFuncCall).Arg0) // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, s2.GetFunc, func(t assert.TestingT, call interface{}) bool {
		return assert.Equal(t, []string{"bar"}, call.(SourcerGetFuncCall).Arg0) // TODO - ergonomics
	})
	mockassert.CalledOnceMatching(t, s3.GetFunc, func(t assert.TestingT, call interface{}) bool {
		return assert.Equal(t, []string{"foo"}, call.(SourcerGetFuncCall).Arg0) // TODO - ergonomics
	})
}

func TestMultiSourcerSkip(t *testing.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.TagsFunc.SetDefaultReturn([]string{"a"})
	s2.TagsFunc.SetDefaultReturn([]string{"b"})
	s3.TagsFunc.SetDefaultReturn([]string{"a"})

	s1.GetFunc.SetDefaultReturn("", FlagSkip, nil)
	s2.GetFunc.SetDefaultReturn("", FlagSkip, nil)
	s3.GetFunc.SetDefaultReturn("", FlagSkip, nil)

	sourcer := NewMultiSourcer(s3, s2, s1)
	assert.Nil(t, sourcer.Init())

	_, flag, err := sourcer.Get([]string{"", ""})
	assert.Nil(t, err)
	assert.Equal(t, FlagSkip, flag)
}

func TestMultiSourcerDump(t *testing.T) {
	s1 := NewMockSourcer()
	s2 := NewMockSourcer()
	s3 := NewMockSourcer()
	s1.DumpFunc.SetDefaultReturn(map[string]string{"a": "foo"})
	s2.DumpFunc.SetDefaultReturn(map[string]string{"b": "bar", "a": "bonk"})
	s3.DumpFunc.SetDefaultReturn(map[string]string{"c": "baz"})

	sourcer := NewMultiSourcer(s3, s2, s1)
	expected := map[string]string{
		"a": "bonk",
		"b": "bar",
		"c": "baz",
	}

	assert.Nil(t, sourcer.Init())
	assert.Equal(t, expected, sourcer.Dump())
}
