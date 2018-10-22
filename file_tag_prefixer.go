package zubrin

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
)

// FileTagPrefixer is a tag modifier which adds a prefix to the values of
// `file` tags. This can be used to register one config multiple times and
// have their initialization be read from different keysl in a config file.
type FileTagPrefixer struct {
	prefix string
}

// NewFileTagPrefixer creates a new FileTagPrefixer.
func NewFileTagPrefixer(prefix string) TagModifier {
	return &FileTagPrefixer{
		prefix: prefix,
	}
}

// AlterFieldTag adds the file prefixer's prefix to the `file` tag value, if one is set.
func (p *FileTagPrefixer) AlterFieldTag(fieldType reflect.StructField, tags *structtag.Tags) error {
	tag, err := tags.Get("file")
	if err != nil {
		return nil
	}

	return tags.Set(&structtag.Tag{
		Key:  "file",
		Name: fmt.Sprintf("%s_%s", p.prefix, tag.Name),
	})
}
