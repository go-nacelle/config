package zubrin

import "github.com/go-nacelle/config/tags"

type TagModifier = tags.TagModifier

var (
	ApplyTagModifiers   = tags.ApplyTagModifiers
	NewEnvTagPrefixer   = tags.NewEnvTagPrefixer
	NewFileTagPrefixer  = tags.NewFileTagPrefixer
	NewDefaultTagSetter = tags.NewDefaultTagSetter
	NewDisplayTagSetter = tags.NewDisplayTagSetter
	NewFileTagSetter    = tags.NewFileTagSetter
)
