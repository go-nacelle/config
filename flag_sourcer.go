package config

import (
	"fmt"
	"strings"
)

type flagSourcer struct {
	args   []string
	parsed map[string]string
}

var _ Sourcer = &flagSourcer{}

// NewFlagSourcer creates a Sourcer that pulls values from the application
// flags.
func NewFlagSourcer(configs ...FlagSourcerConfigFunc) Sourcer {
	return &flagSourcer{args: getFlagSourcerConfigOptions(configs).args}
}

func (s *flagSourcer) Tags() []string {
	return []string{FlagTag}
}

func (s *flagSourcer) Get(values []string) (string, SourcerFlag, error) {
	if values[0] == "" {
		return "", FlagSkip, nil
	}

	if err := s.ensureParsed(); err != nil {
		return "", 0, err
	}

	if val, ok := s.parsed[values[0]]; ok {
		return val, FlagFound, nil
	}

	return "", FlagMissing, nil
}

func (s *flagSourcer) Assets() []string {
	return []string{"<application flags>"}
}

func (s *flagSourcer) Dump() (map[string]string, error) {
	if err := s.ensureParsed(); err != nil {
		return nil, err
	}

	values := map[string]string{}
	for name, value := range s.parsed {
		values[name] = value
	}

	return values, nil
}

func (s *flagSourcer) ensureParsed() error {
	if s.parsed != nil {
		return nil
	}
	s.parsed = map[string]string{}

	for len(s.args) > 0 {
		name, value, err := s.parseOne()
		if err != nil {
			return err
		}

		s.parsed[name] = value
	}

	return nil
}

func (s *flagSourcer) parseOne() (string, string, error) {
	arg := s.args[0]
	s.args = s.args[1:]

	numMinuses := 0
	for i := 0; i < len(arg) && arg[i] == '-'; i++ {
		numMinuses++
	}

	if len(arg) == 0 || numMinuses == 0 || numMinuses > 2 || (len(arg) > numMinuses && arg[numMinuses] == '=') {
		return "", "", fmt.Errorf("illegal flag: %s", arg)
	}

	name := arg[numMinuses:]
	value := ""

	if parts := strings.SplitN(name, "=", 2); len(parts) == 2 {
		name, value = parts[0], parts[1]
	} else {
		if len(s.args) > 0 {
			value, s.args = s.args[0], s.args[1:]
		} else {
			return "", "", fmt.Errorf("flag needs an argument: -%s", name)
		}
	}

	return name, value, nil
}
