package config

import (
	"reflect"
	"testing"

	"github.com/fatih/structtag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ensureEquals(t *testing.T, sourcer Sourcer, values []string, expected string) {
	val, flag, err := sourcer.Get(values)
	require.Nil(t, err)
	assert.Equal(t, FlagFound, flag)
	assert.Equal(t, expected, val)
}

func ensureMatches(t *testing.T, sourcer Sourcer, values []string, expected string) {
	val, flag, err := sourcer.Get(values)
	require.Nil(t, err)
	assert.Equal(t, FlagFound, flag)
	assert.JSONEq(t, expected, val)
}

func ensureMissing(t *testing.T, sourcer Sourcer, values []string) {
	_, flag, err := sourcer.Get(values)
	require.Nil(t, err)
	assert.Equal(t, FlagMissing, flag)
}

func gatherTags(obj interface{}, name string) map[string]string {
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	objType := objValue.Type()

	return gatherTagsStruct(objValue, objType, name)
}

func gatherTagsStruct(objValue reflect.Value, objType reflect.Type, name string) map[string]string {
	for i := 0; i < objType.NumField(); i++ {
		field := objValue.Field(i)
		fieldType := objType.Field(i)

		if fieldType.Anonymous {
			if tags := gatherTagsStruct(field, fieldType.Type, name); tags != nil {
				return tags
			}
		}

		if fieldType.Name == name {
			if tags, ok := getTags(fieldType); ok {
				return decomposeTags(tags)
			}
		}
	}

	return nil
}

func decomposeTags(tags *structtag.Tags) map[string]string {
	fieldTags := map[string]string{}

	for _, name := range tags.Keys() {
		tag, _ := tags.Get(name)
		fieldTags[name] = tag.Name
	}

	return fieldTags
}
