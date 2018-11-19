package zubrin

import (
	"reflect"

	"github.com/fatih/structtag"
)

// DisplayTagSetter is a tag modifier which sets the value of the display
// tag to be the same asn the env tag.
type DisplayTagSetter struct{}

// NewFileTagSetter creates a new DisplayTagSetter.
func NewDisplayTagSetter() TagModifier {
	return &DisplayTagSetter{}
}

// AlterFieldTag sets the value of the display tag to the value fo the env tag.
func (s *DisplayTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tagValue, err := tags.Get("env")
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  "display",
		Name: tagValue.Name,
	})
}
