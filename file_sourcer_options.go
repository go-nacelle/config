package config

type fileSourcerOptions struct {
	fs FileSystem
}

// FileSourcerConfigFunc is a function used to configure instances of
// file sourcers.
type FileSourcerConfigFunc func(*fileSourcerOptions)

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
