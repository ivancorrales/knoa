package internal

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

type outputFormat int32

const (
	YAML outputFormat = iota
	JSON
)

type OutputOpt func(output *Outputter)

func WithOutputFormat(fmt outputFormat) func(output *Outputter) {
	return func(output *Outputter) {
		output.Fmt = fmt
	}
}

type Outputter struct {
	Fmt outputFormat
}

func NewOutput(options ...OutputOpt) *Outputter {
	o := &Outputter{
		Fmt: JSON,
	}
	for _, opt := range options {
		opt(o)
	}
	return o
}

func (o *Outputter) String(content any) string {
	switch o.Fmt {
	case YAML:
		return o.YAML(content)
	case JSON:
		return o.JSON(content)
	}
	return ""
}

func (o *Outputter) YAML(content any) string {
	b, _ := yaml.Marshal(content)
	return string(b)
}

func (o *Outputter) JSON(content any) string {
	b, _ := json.Marshal(content)
	return string(b)
}
