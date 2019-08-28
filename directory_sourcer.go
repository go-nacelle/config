package config

import (
	"path/filepath"
)

// NewOptionalDirectorySourcer creates a directory sourcer if the provided directory
// exists. If the provided file is not found, a sourcer is returned returns no values.
func NewOptionalDirectorySourcer(dirname string, parser FileParser, configs ...DirectorySourcerConfigFunc) Sourcer {
	options := getDirectorySourcerConfigOptions(configs)

	exists, err := options.fs.Exists(dirname)
	if err != nil {
		return newErrorSourcer(err)
	}

	if !exists {
		return &fileSourcer{values: map[string]string{}}
	}

	return NewDirectorySourcer(dirname, parser, configs...)
}

// NewDirectorySourcer creates a sourcer that reads files from a directory. For
// details on parsing format, refer to NewFileParser. Each file in a directory
// is read in alphabetical order. Nested directories are ignored when reading
// directory content, and each found regular file is assumed to be parseable by
// the given FileParser.
func NewDirectorySourcer(dirname string, parser FileParser, configs ...DirectorySourcerConfigFunc) Sourcer {
	options := getDirectorySourcerConfigOptions(configs)

	filenames, err := options.fs.ReadDir(dirname)
	if err != nil {
		return newErrorSourcer(err)
	}

	sourcers := []Sourcer{}
	for _, filename := range filenames {
		sourcers = append(sourcers, NewFileSourcer(filepath.Join(dirname, filename), parser))
	}

	return NewMultiSourcer(sourcers...)
}
