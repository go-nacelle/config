package zubrin

import (
	"reflect"

	"github.com/fatih/structtag"
)

// FileTagSetter is a tag modifier which sets the value of the file
// tag to be the same asn the env tag.
type FileTagSetter struct{}

// NewFileTagSetter creates a new FileTagSetter.
func NewFileTagSetter() TagModifier {
	return &FileTagSetter{}
}

// AlterFieldTag sets the value of the file tag to the value of the env tag.
func (s *FileTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tagValue, err := tags.Get("env")
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  "file",
		Name: tagValue.Name,
	})
}
