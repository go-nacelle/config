package config

type (
	directorySourcerOptions struct{ fs FileSystem }

	// DirectorySourcerConfigFunc is a function used to configure instances of
	// directory sourcers.
	DirectorySourcerConfigFunc func(*directorySourcerOptions)
)

// WithDirectorySourcerFS sets the FileSystem instance.
func WithDirectorySourcerFS(fs FileSystem) DirectorySourcerConfigFunc {
	return func(o *directorySourcerOptions) { o.fs = fs }
}

func getDirectorySourcerConfigOptions(configs []DirectorySourcerConfigFunc) *directorySourcerOptions {
	options := &directorySourcerOptions{
		fs: &realFileSystem{},
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
