package tags

import (
	"reflect"

	"github.com/fatih/structtag"

	"github.com/efritz/zubrin/internal/consts"
)

type fileTagSetter struct{}

// NewFileTagSetter creates a new TagModifier which sets the value
// of the file tag to be the same asn the env tag.
func NewFileTagSetter() TagModifier {
	return &fileTagSetter{}
}

func (s *fileTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tagValue, err := tags.Get(consts.EnvTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  consts.FileTag,
		Name: tagValue.Name,
	})
}
