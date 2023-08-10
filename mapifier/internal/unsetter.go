package internal

type Unsetter struct {
	prefix     string
	funcPrefix func(string) string
}

type UnsetterOpt func(unsetter *Unsetter)

func NewUnsetter(opts ...UnsetterOpt) *Unsetter {
	s := &Unsetter{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Unsetter) Unset(parser *Parser, paths []string) (mutators []Mutator) {
	for _, path := range paths {
		if s.prefix != "" {
			path = s.prefix + path
		}
		if s.funcPrefix != nil {
			path = s.funcPrefix(path)
		}
		m := parser.Parse(path)
		m.operation = UnsetOp
		if m != nil {
			mutators = append(mutators, *m)
		}
	}
	return mutators
}
