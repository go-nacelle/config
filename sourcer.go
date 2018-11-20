package zubrin

import "github.com/efritz/zubrin/sourcer"

type Sourcer = sourcer.Sourcer

var (
	NewMultiSourcer             = sourcer.NewMultiSourcer
	NewGlobSourcer              = sourcer.NewGlobSourcer
	NewEnvSourcer               = sourcer.NewEnvSourcer
	NewFileSourcer              = sourcer.NewFileSourcer
	NewYAMLFileSourcer          = sourcer.NewYAMLFileSourcer
	NewTOMLFileSourcer          = sourcer.NewTOMLFileSourcer
	ParseYAML                   = sourcer.ParseYAML
	ParseTOML                   = sourcer.ParseTOML
	NewOptionalFileSourcer      = sourcer.NewOptionalFileSourcer
	NewDirectorySourcer         = sourcer.NewDirectorySourcer
	NewOptionalDirectorySourcer = sourcer.NewOptionalDirectorySourcer
)
