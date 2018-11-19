package logging

// ConfigFunc is a function used to initialize a logging config.
type ConfigFunc func(*loggingConfig)

// WithMaskedKeys sets the keys which are not displayed when the
// full source content is output to the logger.
func WithMaskedKeys(maskedKeys []string) ConfigFunc {
	return func(l *loggingConfig) { l.maskedKeys = maskedKeys }
}
