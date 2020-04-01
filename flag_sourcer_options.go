package config

import "os"

type (
	flagSourcerOptions struct {
		args []string
	}

	// FlagSourcerConfigFunc is a function used to configure instances of
	// flag sourcers.
	FlagSourcerConfigFunc func(*flagSourcerOptions)
)

// WithFlagSourcerArgs sets raw command line arguments.
func WithFlagSourcerArgs(args []string) FlagSourcerConfigFunc {
	return func(o *flagSourcerOptions) { o.args = args }
}

func getFlagSourcerConfigOptions(configs []FlagSourcerConfigFunc) *flagSourcerOptions {
	options := &flagSourcerOptions{
		args: os.Args[1:],
	}

	for _, f := range configs {
		f(options)
	}

	return options
}
