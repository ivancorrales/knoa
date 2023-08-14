package outputter

import "encoding/json"

type JSON struct {
	pretty bool
	prefix string
	ident  string
}

type JSONOpt func(json *JSON)

func WithPrefixAndIdent(prefix, ident string) func(j *JSON) {
	return func(j *JSON) {
		j.prefix = prefix
		j.ident = ident
		j.pretty = true
	}
}

func (j *JSON) Marshal(content any) (string, error) {
	if j.pretty {
		b, err := json.MarshalIndent(content, j.prefix, j.ident)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	b, err := json.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NewJSON(opts ...JSONOpt) *JSON {
	j := &JSON{}
	for _, opt := range opts {
		opt(j)
	}
	return j
}
