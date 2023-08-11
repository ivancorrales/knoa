package mutator

import (
	"github.com/fatih/structs"
	"github.com/ivancorrales/knoa/sanitizer"
	"reflect"
)

type operationCode int32

const (
	setOp operationCode = iota
	unsetOp
)

type operation struct {
	prefix     string
	funcPrefix func(string) string
}

type OperationOpt func(opt *operation)

func WithFuncPrefix(fn func(string) string) func(o *operation) {
	return func(o *operation) {
		o.funcPrefix = fn
	}
}

func WithStringPrefix(prefix string) func(setter *operation) {
	return func(o *operation) {
		o.prefix = prefix
	}
}

func NewOperation(opts ...OperationOpt) *operation {
	op := &operation{}
	for _, opt := range opts {
		opt(op)
	}
	return op
}

func (op *operation) Set(parser *Parser, pathValueList sanitizer.PathValueList) (mutators []Mutator) {
	for _, pathValue := range pathValueList {
		v := op.checkValue(pathValue.Value)
		path := pathValue.Path
		if op.prefix != "" {
			path = op.prefix + path
		}
		if op.funcPrefix != nil {
			path = op.funcPrefix(path)
		}
		m := parser.Parse(path)
		if m != nil {
			m.addValueToNode(v)
			mutators = append(mutators, *m)
		}
	}
	return mutators
}

func (op *operation) Unset(parser *Parser, paths []string) (mutators []Mutator) {
	for _, path := range paths {
		if op.prefix != "" {
			path = op.prefix + path
		}
		if op.funcPrefix != nil {
			path = op.funcPrefix(path)
		}
		m := parser.Parse(path)
		m.operation = unsetOp
		if m != nil {
			mutators = append(mutators, *m)
		}
	}
	return mutators
}

func (op *operation) checkValue(value any) any {
	switch reflect.ValueOf(value).Kind() {
	case reflect.Struct:
		mapValues := structs.Map(value)
		return mapValues
	default:
		return value
	}
}