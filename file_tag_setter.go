package config

import (
	"reflect"

	"github.com/fatih/structtag"
)

type fileTagSetter struct{}

var _ TagModifier = &fileTagSetter{}

// NewFileTagSetter creates a new TagModifier which sets the value
// of the file tag to be the same as the env tag.
func NewFileTagSetter() TagModifier {
	return &fileTagSetter{}
}

func (s *fileTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tagValue, err := tags.Get(EnvTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  FileTag,
		Name: tagValue.Name,
	})
}
