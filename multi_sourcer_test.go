package config

import (
	"sort"
	"testing"

	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.Nil(t, sourcer.Init())

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
	require.Nil(t, sourcer.Init())
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
	require.Nil(t, sourcer.Init())

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
	require.Nil(t, sourcer.Init())

	_, flag, err := sourcer.Get([]string{"foo", "bar"})
	require.Nil(t, err)
	assert.Equal(t, FlagMissing, flag)

	mockassert.CalledOnceWith(t, s1.GetFunc, mockassert.Values([]string{"foo"}))
	mockassert.CalledOnceWith(t, s2.GetFunc, mockassert.Values([]string{"bar"}))
	mockassert.CalledOnceWith(t, s3.GetFunc, mockassert.Values([]string{"foo"}))
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
	require.Nil(t, sourcer.Init())

	_, flag, err := sourcer.Get([]string{"", ""})
	require.Nil(t, err)
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

	require.Nil(t, sourcer.Init())
	assert.Equal(t, expected, sourcer.Dump())
}
