package config

type errorSourcer struct {
	err error
}

func newErrorSourcer(err error) Sourcer {
	return &errorSourcer{err}
}

func (s *errorSourcer) Tags() []string {
	return nil
}

func (s *errorSourcer) Get(values []string) (string, SourcerFlag, error) {
	return "", FlagUnknown, s.err
}

func (s *errorSourcer) Assets() []string {
	return nil
}

func (s *errorSourcer) Dump() map[string]string {
	return nil
}
