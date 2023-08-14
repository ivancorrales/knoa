package knoa

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"

	"github.com/ivancorrales/knoa/internal"
	"github.com/ivancorrales/knoa/mutator"
	"github.com/ivancorrales/knoa/outputter"
	"github.com/ivancorrales/knoa/sanitizer"
)

func FromMap(content map[string]any, opts ...Opt) Knoa[map[string]any] {
	return load[map[string]any](content, opts...)
}

func Map(opts ...Opt) Knoa[map[string]any] {
	return New[map[string]any](opts...)
}

func FromArray(content []any, opts ...Opt) Knoa[[]any] {
	return load[[]any](content, opts...)
}

func Array(opts ...Opt) Knoa[[]any] {
	return New[[]any](opts...)
}

func load[T Type](content T, options ...Opt) Knoa[T] {
	b := &builder{
		strictMode:  false,
		attrNameFmt: mutator.DefAttributeNameFormat,
	}
	for _, opt := range options {
		opt(b)
	}
	pathRegExp, attrRegExpr := mutator.RegExpsFromAttributeFormat(b.attrNameFmt)
	c, _ := internal.Normalize(content).(T)
	return &knoa[T]{
		strictMode: b.strictMode,
		parser: &mutator.Parser{
			Strict:          b.strictMode,
			RegExp:          pathRegExp,
			AttributeRegExp: attrRegExpr,
		},
		content: c,
	}
}

type Type internal.Type

type Knoa[T Type] interface {
	Set(pathValueList ...any) Knoa[T]
	Unset(pathValueList ...string) Knoa[T]
	Apply(args ...any) Knoa[T]
	With(opts ...mutator.OperationOpt) func(pathValueList ...any) Knoa[T]
	Out() T
	YAML(opts ...outputter.YAMLOpt) string
	JSON(opts ...outputter.JSONOpt) string
	To(output interface{})
	Error() error
}

type knoa[T Type] struct {
	strictMode bool
	mutators   []mutator.Mutator
	parser     *mutator.Parser
	content    T
	err        error
}

type Opt func(sanitizer *builder)

type builder struct {
	strictMode  bool
	attrNameFmt string
}

func WithStrictMode(strict bool) func(builder *builder) {
	return func(builder *builder) {
		builder.strictMode = strict
	}
}

func WithAttributeNameFormat(attrNameFmt string) func(builder *builder) {
	return func(opts *builder) {
		opts.attrNameFmt = attrNameFmt
	}
}

func New[T Type](options ...Opt) Knoa[T] {
	var content T
	return load[T](content, options...)
}

func (k *knoa[T]) With(opts ...mutator.OperationOpt) func(args ...any) Knoa[T] {
	setter := mutator.NewOperation(opts...)
	return func(args ...any) Knoa[T] {
		pathValueList := sanitizer.SanitizePathValueList(k.strictMode, args...)
		k.mutators = append(k.mutators, setter.Set(k.parser, pathValueList)...)
		return k
	}
}

func (k *knoa[T]) Set(args ...any) Knoa[T] {
	pathValueList := sanitizer.SanitizePathValueList(k.strictMode, args...)
	k.mutators = append(k.mutators, mutator.NewOperation().Set(k.parser, pathValueList)...)
	return k
}

func (k *knoa[T]) Unset(args ...string) Knoa[T] {
	k.mutators = append(k.mutators, mutator.NewOperation().Unset(k.parser, args)...)
	return k
}

func (k *knoa[T]) Apply(args ...any) Knoa[T] {
	pathFuncList := sanitizer.SanitizePathFuncList(k.strictMode, args...)
	k.mutators = append(k.mutators, mutator.NewOperation().Apply(k.parser, pathFuncList)...)
	return k
}

func (k *knoa[T]) Out() T {
	var content T = k.content
	for _, m := range k.mutators {
		switch reflect.ValueOf(content).Kind() {
		case reflect.Slice, reflect.Array:
			in, ok := reflect.ValueOf(content).Interface().([]any)
			if !ok {
				k.err = fmt.Errorf("unsupported array type")
				break
			}
			arrayIn, err := m.Child().ToArray(in)
			if err != nil {
				k.err = err
				break
			}
			content, _ = reflect.ValueOf(arrayIn).Interface().(T)
		case reflect.Map:
			in, ok := reflect.ValueOf(content).Interface().(map[string]any)
			if !ok {
				k.err = fmt.Errorf("unsupported map type")
				break
			}
			mapIn, err := m.Child().ToMap(in)
			if err != nil {
				k.err = err
				break
			}
			content, _ = reflect.ValueOf(mapIn).Interface().(T)
		default:
			k.err = fmt.Errorf("unsupporteed output type '%s'", reflect.TypeOf(content).Kind())
		}
	}
	return content
}

func (k *knoa[T]) YAML(opts ...outputter.YAMLOpt) string {
	content := k.Out()
	str, err := outputter.NewYAML(opts...).Marshal(content)
	k.err = err
	return str
}

func (k *knoa[T]) JSON(opts ...outputter.JSONOpt) string {
	content := k.Out()
	str, err := outputter.NewJSON(opts...).Marshal(content)
	k.err = err
	return str
}

func (k *knoa[T]) To(out interface{}) {
	content := k.Out()
	k.err = mapstructure.Decode(content, out)
}

func (k *knoa[T]) Error() error {
	return k.err
}
