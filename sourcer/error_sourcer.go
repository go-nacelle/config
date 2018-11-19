package sourcer

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
