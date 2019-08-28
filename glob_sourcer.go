package config

// NewGlobSourcer creates a sourcer that reads all files that match
// the given glob pattern. Each matching file is read in alphabetical
// order of path. Each matching pathis assumed to be parseable by the
// given FileParser.
func NewGlobSourcer(pattern string, parser FileParser, configs ...GlobSourcerConfigFunc) Sourcer {
	options := getGlobSourcerConfigOptions(configs)

	sourcers := []Sourcer{}
	if paths, err := options.fs.Glob(pattern); err == nil { // TODO - configure
		for _, path := range paths {
			sourcers = append(sourcers, NewFileSourcer(path, parser))
		}
	}

	return NewMultiSourcer(sourcers...)
}
