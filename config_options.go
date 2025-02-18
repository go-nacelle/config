package config

type configOptions struct {
	logger          Logger
	maskedKeys      []string
	sourceDumpOnErr bool
}

// ConfigOptionsFunc is a function used to configure instances of config.
type ConfigOptionsFunc func(*configOptions)

// WithLogger sets the Logger instance.
func WithLogger(logger Logger) ConfigOptionsFunc {
	return func(o *configOptions) { o.logger = logger }
}

// WithMaskedKeys sets the field names of values masked in log messages.
func WithMaskedKeys(maskedKeys []string) ConfigOptionsFunc {
	return func(o *configOptions) { o.maskedKeys = maskedKeys }
}

// WithSourceDumpOnError sets whether the config source will be dumped on a
// config load error or not. It can be useful for troubleshooting, but unless
// other options are set correctly it could lead to secrets being logged.
func WithSourceDumpOnError(sourceDumpOnErr bool) ConfigOptionsFunc {
	return func(o *configOptions) { o.sourceDumpOnErr = sourceDumpOnErr }
}

func getConfigOptions(configs []ConfigOptionsFunc) *configOptions {
	options := &configOptions{
		logger: &nilLogger{},
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
