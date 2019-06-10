package config

import (
	"reflect"

	"github.com/fatih/structtag"
)

type defaultTagSetter struct {
	field        string
	defaultValue string
}

// NewDefaultTagSetter creates a new TagModifier which sets the value of
// the default tag for a particular field. This is used to change the default
// values provided by third party libraries (for which a source change would
// be otherwise required).
func NewDefaultTagSetter(field string, defaultValue string) TagModifier {
	return &defaultTagSetter{
		field:        field,
		defaultValue: defaultValue,
	}
}

func (s *defaultTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	if fieldType.Name != s.field {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  DefaultTag,
		Name: s.defaultValue,
	})
}
