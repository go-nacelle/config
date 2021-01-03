package config

type configOptions struct {
	logger     Logger
	maskedKeys []string
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

func getConfigOptions(configs []ConfigOptionsFunc) *configOptions {
	options := &configOptions{
		logger: &nilLogger{},
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
