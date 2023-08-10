package internal

import (
	"reflect"

	"github.com/fatih/structs"
)

type Setter struct {
	prefix     string
	funcPrefix func(string) string
}

type SetterOpt func(setter *Setter)

func WithFuncPrefix(fn func(string) string) func(setter *Setter) {
	return func(setter *Setter) {
		setter.funcPrefix = fn
	}
}

func WithStringPrefix(prefix string) func(setter *Setter) {
	return func(setter *Setter) {
		setter.prefix = prefix
	}
}

func NewSetter(opts ...SetterOpt) *Setter {
	s := &Setter{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Setter) Set(parser *Parser, pathValueList PathValueList) (mutators []Mutator) {
	for _, pathValue := range pathValueList {
		v := s.checkValue(pathValue.Value)
		path := pathValue.Path
		if s.prefix != "" {
			path = s.prefix + path
		}
		if s.funcPrefix != nil {
			path = s.funcPrefix(path)
		}
		m := parser.Parse(path)
		if m != nil {
			m.WithValue(v)
			mutators = append(mutators, *m)
		}
	}
	return mutators
}

func (s *Setter) checkValue(value any) any {
	switch reflect.ValueOf(value).Kind() {
	case reflect.Struct:
		mapValues := structs.Map(value)
		return mapValues
	default:
		return value
	}
}
