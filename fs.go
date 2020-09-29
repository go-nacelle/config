package config

import (
	"io/ioutil"
	"os"

	"github.com/mattn/go-zglob"
)

// FileSystem is an interface wrapping filesystem operations for
// sourcers that read information from disk. This is necessary in
// order to allow remote and in-memory filesystems that may be
// present in some application deployments. Third-party libraries
// such as  spf13/afero that provide FS-like functionality can be
// shimmed into this interface.
type FileSystem interface {
	// Exists determines if the given path exists.
	Exists(path string) (bool, error)

	// ListFiles returns the names of the files that are a direct
	// child of the directory at the given path.
	ListFiles(path string) ([]string, error)

	// Glob returns the paths that the given pattern matches.
	Glob(pattern string) ([]string, error)

	// ReadFile returns the content of the file at the given path.
	ReadFile(path string) ([]byte, error)
}

type realFileSystem struct{}

var _ FileSystem = &realFileSystem{}

func (fs *realFileSystem) Exists(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return false, err
	}

	return true, nil
}

func (fs *realFileSystem) ListFiles(path string) ([]string, error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fileEntries := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileEntries = append(fileEntries, entry.Name())
	}

	return fileEntries, nil
}

func (fs *realFileSystem) Glob(pattern string) ([]string, error) {
	return zglob.Glob(pattern)
}

func (fs *realFileSystem) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
