package config

type (
	globSourcerOptions struct{ fs FileSystem }

	// GlobSourcerConfigFunc is a function used to configure instances of
	// glob sourcers.
	GlobSourcerConfigFunc func(*globSourcerOptions)
)

// WithGlobSourcerFS sets the FileSystem instance.
func WithGlobSourcerFS(fs FileSystem) GlobSourcerConfigFunc {
	return func(o *globSourcerOptions) { o.fs = fs }
}

func getGlobSourcerConfigOptions(configs []GlobSourcerConfigFunc) *globSourcerOptions {
	options := &globSourcerOptions{
		fs: &realFileSystem{},
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
