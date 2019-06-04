package tags

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"

	"github.com/go-nacelle/config/internal/consts"
)

type fileTagPrefixer struct {
	prefix string
}

// NewFileTagPrefixer creates a new TagModifier which adds a prefix to the
// values of `file` tags. This can be used to register one config multiple
// times and have their initialization be read from different keysl in a
// config file.
func NewFileTagPrefixer(prefix string) TagModifier {
	return &fileTagPrefixer{
		prefix: prefix,
	}
}

func (p *fileTagPrefixer) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tag, err := tags.Get(consts.FileTag)
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  consts.FileTag,
		Name: fmt.Sprintf("%s_%s", p.prefix, tag.Name),
	})
}
