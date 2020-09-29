package config

// Sourcer pulls requested names from a variable source. This can be the
// environment, a file, a remote server, etc. This can be done on-demand
// per variable, or a cache of variables can be built on startup and then
// pulled from a cached mapping as requested.
type Sourcer interface {
	// Init is a hook for certain classes of sourcers to read and normalize
	// the source data. This gives a canonical place where external errors
	// can occur that are not directly related to validation.
	Init() error

	// Tags returns a list of tags which are required to get a value from
	// the source. Order matters.
	Tags() []string

	// Get will retrieve a value from the source with the given tag values.
	// The tag values passed to this method will be in the same order as
	// returned from the Tags method. The flag return value directs config
	// population whether or not this value should be treated as missing or
	// skippable.
	Get(values []string) (string, SourcerFlag, error)

	// Assets returns a list of names of assets that compose the sourcer.
	// This can be a list of matched files that are read, or a token that
	// denotes a fixed source.
	Assets() []string

	// Dump returns the full content of the sourcer. This is used by the
	// logging package to show the content of the environment and config
	// files when a value is missing or otherwise illegal.
	Dump() map[string]string
}

type SourcerFlag int

const (
	FlagUnknown SourcerFlag = iota
	FlagFound
	FlagMissing
	FlagSkip
)
