package zubrin

import (
	"github.com/mattn/go-zglob"
)

// NewGlobSourcer creates a sourcer that reads all files that match
// the given glob pattern. Each matching file is read in alphabetical
// order of path. Each matching pathis assumed to be parseable by the
// given FileParser.
func NewGlobSourcer(pattern string, parser FileParser) Sourcer {
	sourcers := []Sourcer{}
	if paths, err := zglob.Glob(pattern); err == nil {
		for _, path := range paths {
			sourcers = append(sourcers, NewFileSourcer(path, parser))
		}
	}

	return NewMultiSourcer(sourcers...)
}
