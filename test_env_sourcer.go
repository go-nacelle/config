package config

import (
	"strings"
)

type testEnvSourcer struct {
	values map[string]string
}

var _ Sourcer = &testEnvSourcer{}

// NewTestEnvSourcer creates a Sourcer that returns values from a given
// map.
func NewTestEnvSourcer(values map[string]string) Sourcer {
	normalized := map[string]string{}
	for key, value := range values {
		normalized[strings.ToUpper(key)] = value
	}

	return &testEnvSourcer{
		values: normalized,
	}
}

func (s *testEnvSourcer) Init() error {
	return nil
}

func (s *testEnvSourcer) Tags() []string {
	return []string{EnvTag}
}

func (s *testEnvSourcer) Get(values []string) (string, SourcerFlag, error) {
	if values[0] == "" {
		return "", FlagSkip, nil
	}

	if val, ok := s.values[strings.ToUpper(values[0])]; ok {
		return val, FlagFound, nil
	}

	return "", FlagMissing, nil
}

func (s *testEnvSourcer) Assets() []string {
	return []string{"<test environment>"}
}

func (s *testEnvSourcer) Dump() map[string]string {
	return s.values
}
