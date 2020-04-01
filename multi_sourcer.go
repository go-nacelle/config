package config

import (
	"sort"
)

type multiSourcer struct {
	sourcers []Sourcer
	tags     []string
}

var _ Sourcer = &multiSourcer{}

// NewMultiSourcer creates a sourcer that reads form each sourcer.
// The last value found is returned - sourcers should be provided
// from low priority to high priority.
func NewMultiSourcer(sourcers ...Sourcer) Sourcer {
	set := map[string]struct{}{}
	for _, sourcer := range sourcers {
		for _, tag := range sourcer.Tags() {
			set[tag] = struct{}{}
		}
	}

	tags := []string{}
	for tag := range set {
		tags = append(tags, tag)
	}

	sort.Strings(tags)

	return &multiSourcer{
		sourcers: sourcers,
		tags:     tags,
	}
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

func (s *multiSourcer) Dump() (map[string]string, error) {
	values := map[string]string{}
	for i := len(s.sourcers) - 1; i >= 0; i-- {
		dump, err := s.sourcers[i].Dump()
		if err != nil {
			return nil, err
		}

		for k, v := range dump {
			values[k] = v
		}
	}

	return values, nil
}
