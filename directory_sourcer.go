package config

import (
	"path/filepath"
)

type directorySourcer struct {
	dirname  string
	parser   FileParser
	fs       FileSystem
	optional bool
	sourcer  Sourcer
}

var _ Sourcer = &directorySourcer{}

// NewOptionalDirectorySourcer creates a directory sourcer that does not error on init
// if the directory does not exist. In this case, the underlying sourcer returns no
// values.
func NewOptionalDirectorySourcer(dirname string, parser FileParser, configs ...DirectorySourcerConfigFunc) Sourcer {
	return &directorySourcer{
		dirname:  dirname,
		parser:   parser,
		fs:       getDirectorySourcerConfigOptions(configs).fs,
		optional: true,
	}
}

// NewDirectorySourcer creates a sourcer that reads files from a directory. For
// details on parsing format, refer to NewFileParser. Each file in a directory
// is read in alphabetical order. Nested directories are ignored when reading
// directory content, and each found regular file is assumed to be parseable by
// the given FileParser.
func NewDirectorySourcer(dirname string, parser FileParser, configs ...DirectorySourcerConfigFunc) Sourcer {
	return &directorySourcer{
		dirname: dirname,
		parser:  parser,
		fs:      getDirectorySourcerConfigOptions(configs).fs,
	}
}

func (s *directorySourcer) Init() error {
	sourcer, err := s.getSourcer()
	if err != nil {
		return err
	}

	if err := sourcer.Init(); err != nil {
		return err
	}

	s.sourcer = sourcer
	return nil
}

func (s *directorySourcer) getSourcer() (Sourcer, error) {
	if s.optional {
		exists, err := s.fs.Exists(s.dirname)
		if err != nil {
			return nil, err
		}

		if !exists {
			return NewTestEnvSourcer(nil), nil
		}
	}

	filenames, err := s.fs.ListFiles(s.dirname)
	if err != nil {
		return nil, err
	}

	sourcers := []Sourcer{}
	for _, filename := range filenames {
		sourcers = append(sourcers, NewFileSourcer(
			filepath.Join(s.dirname, filename),
			s.parser,
			WithFileSourcerFS(s.fs),
		))
	}

	return NewMultiSourcer(sourcers...), nil
}

func (s *directorySourcer) Tags() []string {
	return s.sourcer.Tags()
}

func (s *directorySourcer) Get(values []string) (string, SourcerFlag, error) {
	return s.sourcer.Get(values)
}

func (s *directorySourcer) Assets() []string {
	return s.sourcer.Assets()
}

func (s *directorySourcer) Dump() map[string]string {
	return s.sourcer.Dump()
}
