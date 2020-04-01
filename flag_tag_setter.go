package config

import (
	"reflect"

	"github.com/fatih/structtag"
)

type flagTagSetter struct{}

// NewFlagTagSetter creates a new TagModifier which sets the value
// of the flag tag to be the same as the env tag.
func NewFlagTagSetter() TagModifier {
	return &flagTagSetter{}
}

func (s *flagTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tagValue, err := tags.Get(EnvTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  FlagTag,
		Name: tagValue.Name,
	})
}
