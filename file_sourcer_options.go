package config

type (
	fileSourcerOptions struct{ fs FileSystem }

	// FileSourcerConfigFunc is a functin used to configure instances of
	// file sourcers.
	FileSourcerConfigFunc func(*fileSourcerOptions)
)

// WithFileSourcerFS sets the FileSystem instance.
func WithFileSourcerFS(fs FileSystem) FileSourcerConfigFunc {
	return func(o *fileSourcerOptions) { o.fs = fs }
}

func getFileSourcerConfigOptions(configs []FileSourcerConfigFunc) *fileSourcerOptions {
	options := &fileSourcerOptions{
		fs: &realFileSystem{},
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
