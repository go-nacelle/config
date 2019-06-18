package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// NewOptionalDirectorySourcer creates a directory sourcer if the provided directoy
// exists. the provided file is not found, a sourcer is returned returns no values.
func NewOptionalDirectorySourcer(dirname string, parser FileParser) Sourcer {
	if _, err := os.Stat(dirname); err != nil && os.IsNotExist(err) {
		return &fileSourcer{values: map[string]string{}}
	}

	return NewDirectorySourcer(dirname, parser)
}

// NewDirectorySourcer creates a sourcer that reads files from a directory. For
// details on parsing format, refer to NewFileParser. Each file in a directory
// is read in alphabetical order. Nested directories are ignored when reading
// directory content, and each found regular file is assumed to be parseable by
// the given FileParser.
func NewDirectorySourcer(dirname string, parser FileParser) Sourcer {
	entries, err := ioutil.ReadDir(dirname)
	if err != nil {
		return newErrorSourcer(err)
	}

	sourcers := []Sourcer{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		sourcers = append(sourcers, NewFileSourcer(filepath.Join(dirname, entry.Name()), parser))
	}

	return NewMultiSourcer(sourcers...)
}
