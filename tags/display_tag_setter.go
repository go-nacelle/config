package tags

import (
	"reflect"

	"github.com/fatih/structtag"

	"github.com/efritz/zubrin/internal/consts"
)

type displayTagSetter struct{}

// NewDisplayTagSetter creates a new TagModifier which sets the value
// of the display tag to be the same asn the env tag..
func NewDisplayTagSetter() TagModifier {
	return &displayTagSetter{}
}

func (s *displayTagSetter) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tagValue, err := tags.Get(consts.EnvTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  consts.DisplayTag,
		Name: tagValue.Name,
	})
}
