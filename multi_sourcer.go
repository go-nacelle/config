package config

import (
	"sort"
)

type multiSourcer struct {
	sourcers []Sourcer
	tags     []string
}

var _ Sourcer = &multiSourcer{}

// NewMultiSourcer creates a sourcer that reads form each sourcer. The last value
// found is returned - sourcers should be provided from low priority to high priority.
func NewMultiSourcer(sourcers ...Sourcer) Sourcer {
	return &multiSourcer{sourcers: sourcers}
}

func (s *multiSourcer) Init() error {
	tagSet := map[string]struct{}{}
	for _, sourcer := range s.sourcers {
		if err := sourcer.Init(); err != nil {
			return err
		}

		for _, tag := range sourcer.Tags() {
			tagSet[tag] = struct{}{}
		}
	}

	tags := []string{}
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	sort.Strings(tags)
	s.tags = tags

	return nil
}

func (s *multiSourcer) Tags() []string {
	return s.tags
}

func (s *multiSourcer) Get(values []string) (string, SourcerFlag, error) {
	correlation := map[string]string{}
	for i, value := range values {
		correlation[s.tags[i]] = value
	}

	skip := true
	for i := len(s.sourcers) - 1; i >= 0; i-- {
		sourcer := s.sourcers[i]

		sourcerValues := []string{}
		for _, tag := range sourcer.Tags() {
			sourcerValues = append(sourcerValues, correlation[tag])
		}

		val, flag, err := sourcer.Get(sourcerValues)
		if err != nil {
			return "", FlagMissing, err
		}

		if flag == FlagFound {
			return val, FlagFound, nil
		}

		skip = skip && flag == FlagSkip
	}

	if skip {
		return "", FlagSkip, nil
	}

	return "", FlagMissing, nil
}

func (s *multiSourcer) Assets() []string {
	assets := []string{}
	for _, sourcer := range s.sourcers {
		assets = append(assets, sourcer.Assets()...)
	}

	return assets
}

func (s *multiSourcer) Dump() map[string]string {
	values := map[string]string{}
	for i := len(s.sourcers) - 1; i >= 0; i-- {
		for k, v := range s.sourcers[i].Dump() {
			values[k] = v
		}
	}

	return values
}
