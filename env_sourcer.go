package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type envSourcer struct {
	prefix string
}

var _ Sourcer = &envSourcer{}

var replacePattern = regexp.MustCompile(`[^A-Za-z0-9_]+`)

// NewEnvSourcer creates a Sourcer that pulls values from the environment.
// The environment variable {PREFIX}_{NAME} is read before and, if empty,
// the environment variable {NAME} is read as a fallback. The prefix is
// normalized by replacing all non-alpha characters with an underscore,
// removing leading and trailing underscores, and collapsing consecutive
// underscores with a single character.
func NewEnvSourcer(prefix string) Sourcer {
	prefix = strings.Trim(
		string(replacePattern.ReplaceAll(
			[]byte(prefix),
			[]byte("_"),
		)),
		"_",
	)

	return &envSourcer{prefix: prefix}
}

func (s *envSourcer) Tags() []string {
	return []string{EnvTag}
}

func (s *envSourcer) Get(values []string) (string, SourcerFlag, error) {
	if values[0] == "" {
		return "", FlagSkip, nil
	}

	envvars := []string{
		strings.ToUpper(fmt.Sprintf("%s_%s", s.prefix, values[0])),
		strings.ToUpper(values[0]),
	}

	for _, envvar := range envvars {
		if val, ok := os.LookupEnv(envvar); ok {
			return val, FlagFound, nil
		}
	}

	return "", FlagMissing, nil
}

func (s *envSourcer) Assets() []string {
	return []string{"<os environment>"}
}

func (s *envSourcer) Dump() (map[string]string, error) {
	values := map[string]string{}
	for _, name := range getNames() {
		values[name] = os.Getenv(name)
	}

	return values, nil
}

func getNames() []string {
	names := []string{}
	for _, pair := range os.Environ() {
		names = append(names, strings.Split(pair, "=")[0])
	}

	return names
}
