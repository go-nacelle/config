package config

import (
	"reflect"

	"github.com/fatih/structtag"
)

type displayTagSetter struct{}

// NewDisplayTagSetter creates a new TagModifier which sets the value
// of the display tag to be the same asn the env tag..
func NewDisplayTagSetter() TagModifier {
	return &displayTagSetter{}
}

func (s *displayTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tagValue, err := tags.Get(EnvTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  DisplayTag,
		Name: tagValue.Name,
	})
}
