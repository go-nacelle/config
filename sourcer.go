package zubrin

import "github.com/efritz/zubrin/sourcer"

type Sourcer = sourcer.Sourcer

var (
	NewMultiSourcer             = sourcer.NewMultiSourcer
	NewGlobSourcer              = sourcer.NewGlobSourcer
	NewEnvSourcer               = sourcer.NewEnvSourcer
	NewFileSourcer              = sourcer.NewFileSourcer
	NewOptionalFileSourcer      = sourcer.NewOptionalFileSourcer
	NewDirectorySourcer         = sourcer.NewDirectorySourcer
	NewOptionalDirectorySourcer = sourcer.NewOptionalDirectorySourcer
)
