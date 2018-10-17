package zubrin

import (
	"io/ioutil"
	"path/filepath"
)

type directorySourcer struct {
	err error
}

// NewDirectorySourcer creates a sourcer that reads files from a directory. For
// details on parsing format, refer to NewFileParser. Each file in a directory
// is read in alphabetical order. The directory is assumed to have only files
// and each file must be parseable by the given FileParser.
func NewDirectorySourcer(dirname string, parser FileParser) Sourcer {
	entries, err := ioutil.ReadDir(dirname)
	if err != nil {
		return &fileSourcer{err: err}
	}

	sourcers := []Sourcer{}
	for _, entry := range entries {
		sourcers = append(sourcers, NewFileSourcer(filepath.Join(dirname, entry.Name()), parser))
	}

	return NewMultiSourcer(sourcers...)
}
