package config

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
)

type flagTagPrefixer struct {
	prefix string
}

var _ TagModifier = &flagTagPrefixer{}

// NewFlagTagPrefixer creates a new TagModifier which adds a prefix to the
// values of `flag` tags.
func NewFlagTagPrefixer(prefix string) TagModifier {
	return &flagTagPrefixer{
		prefix: prefix,
	}
}

func (p *flagTagPrefixer) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tag, err := tags.Get(FlagTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  FlagTag,
		Name: fmt.Sprintf("%s_%s", p.prefix, tag.Name),
	})
}
