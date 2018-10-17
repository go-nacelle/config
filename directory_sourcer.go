package zubrin

import (
	"io/ioutil"
	"path/filepath"
)

type directorySourcer struct {
	err error
}

// TODO
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
