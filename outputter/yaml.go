package outputter

import (
	"gopkg.in/yaml.v3"
)

type YAML struct{}

type YAMLOpt func(y *YAML)

func (y *YAML) Marshal(content any) (string, error) {
	b, err := yaml.Marshal(content)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func NewYAML(opts ...YAMLOpt) *YAML {
	y := &YAML{}
	for _, opt := range opts {
		opt(y)
	}
	return y
}
