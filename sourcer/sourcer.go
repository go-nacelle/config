package sourcer

type (
	// Sourcer pulls requested names from a variable source. This can be the
	// environment, a file, a remote server, etc. This can be done on-demand
	// per variable, or a cache of variables can be built on startup and then
	// pulled from a cached mapping as requested.
	Sourcer interface {
		// Tags returns a list of tags which are required to get a value from
		// the source. Order matters.
		Tags() []string

		// Get will retrieve a value from the source with the given tag values.
		// The tag values passed to this method will be in the same order as
		// returned from the Tags method. The flag return value directs config
		// population whether or not this value should be treated as missing or
		// skippable.
		Get(values []string) (string, SourcerFlag, error)
	}

	SourcerFlag int
)

const (
	FlagUnknown SourcerFlag = iota
	FlagFound
	FlagMissing
	FlagSkip
)
