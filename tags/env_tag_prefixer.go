package tags

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"

	"github.com/efritz/zubrin/internal/consts"
)

type envTagPrefixer struct {
	prefix string
}

// NewEnvTagPrefixer creates a new TagModifier which adds a prefix to the
// values of `env` tags. This can be used to register one config multiple
// times and have their initialization be read from different environment
// variables.
func NewEnvTagPrefixer(prefix string) TagModifier {
	return &envTagPrefixer{
		prefix: prefix,
	}
}

func (p *envTagPrefixer) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tag, err := tags.Get(consts.EnvTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  consts.EnvTag,
		Name: fmt.Sprintf("%s_%s", p.prefix, tag.Name),
	})
}
