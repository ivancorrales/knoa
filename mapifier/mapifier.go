package mapifier

import (
	"log"
	"reflect"

	"github.com/fatih/structs"

	"github.com/ivancorrales/knoa/mapifier/internal"
)

type Mapifier[S Type] interface {
	Set(pathValueList ...any) Mapifier[S]
	With(opts ...SetterOpt) func(pathValueList ...any) Mapifier[S]
	Out() S
	YAML() string
	JSON() string
	String(opts ...internal.OutputOpt) string
}

type Type interface {
	map[string]any | []any | any
}

type mapifier[T Type] struct {
	mutators  []internal.Mutator
	sanitizer *internal.Sanitizer
	parser    *internal.Parser
	setter    *internal.Setter
	content   T
}

type KnoaOpt func(sanitizer *builder)

type builder struct {
	strictMode  bool
	attrNameFmt string
}

func New[T Type](options ...KnoaOpt) Mapifier[T] {
	var content T
	return Load[T](content, options...)
}

func convert[T Type](input T) any {
	value := reflect.ValueOf(input)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		itemsLen := value.Len()
		output := make([]any, itemsLen)
		for i := 0; i < itemsLen; i++ {
			itemValue := reflect.ValueOf(input).Index(i)
			switch itemValue.Kind() {
			case reflect.Slice, reflect.Array:
				output[i] = convert(itemValue.Interface().([]any))
			case reflect.Map:
				output[i] = convert(itemValue.Interface().(map[string]any))
			case reflect.Struct:
				v := structs.Map(itemValue.Interface())
				output[i] = convert(v)
			default:
				if reflect.ValueOf(itemValue.Interface()).Kind() == reflect.Struct {
					v := structs.Map(itemValue.Interface())
					output[i] = convert(v)
				} else {
					output[i] = itemValue.Interface()
				}
			}
		}
		if itemsLen == 0 {
			output = make([]any, 1)
		}
		return output
	case reflect.Struct:
		return structs.Map(input)
	case reflect.Map:
		output := make(map[string]any)
		in, ok := reflect.ValueOf(input).Interface().(map[string]any)
		if ok {
			for k, v := range in {
				switch reflect.ValueOf(v).Kind() {
				case reflect.Slice, reflect.Array:
					output[k] = convert(v)
				case reflect.Map:
					output[k] = convert(v.(map[string]any))
				case reflect.Struct:
					output[k] = convert(structs.Map(v))
				default:
					output[k] = v
				}
			}
		}
		return output
	}
	return input
}

func Load[T Type](content T, options ...KnoaOpt) Mapifier[T] {
	b := &builder{
		strictMode:  false,
		attrNameFmt: internal.DefAttributeNameFormat,
	}
	for _, opt := range options {
		opt(b)
	}
	pathRegExp, attrRegExpr := internal.RegExpsFromAttributeFormat(b.attrNameFmt)
	var contentMap T
	procContent := convert(content)
	switch reflect.ValueOf(procContent).Kind() {
	case reflect.Map:
		contentMap = procContent.(T)
	case reflect.Slice, reflect.Array:
		contentMap = procContent.(T)
	}

	p := &mapifier[T]{
		sanitizer: &internal.Sanitizer{
			Strict: b.strictMode,
		},
		parser: &internal.Parser{
			Strict:          b.strictMode,
			RegExp:          pathRegExp,
			AttributeRegExp: attrRegExpr,
		},
		content: contentMap,
	}

	return p
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

func (p *mapifier[T]) With(opts ...SetterOpt) func(args ...any) Mapifier[T] {
	setter := internal.NewSetter(opts...)
	return func(args ...any) Mapifier[T] {
		pathValueList := p.sanitizer.SanitizePathValueList(args...)
		p.mutators = append(p.mutators, setter.Set(p.parser, pathValueList)...)
		return p
	}
}

func (p *mapifier[T]) Set(args ...any) Mapifier[T] {
	pathValueList := p.sanitizer.SanitizePathValueList(args...)
	p.mutators = append(p.mutators, internal.NewSetter().Set(p.parser, pathValueList)...)
	return p
}

func (p *mapifier[T]) Out() T {
	var content T = p.content
	for _, m := range p.mutators {
		switch reflect.ValueOf(content).Kind() {
		case reflect.Slice, reflect.Array:
			in, ok := reflect.ValueOf(content).Interface().([]any)
			if ok {
				content, _ = reflect.ValueOf(m.Child().ToArray(in)).Interface().(T)
			}
		case reflect.Map:
			in, ok := reflect.ValueOf(content).Interface().(map[string]any)
			if ok {
				content, _ = reflect.ValueOf(m.Child().ToMap(in)).Interface().(T)
			}
		default:
			log.Fatalf("unsupporteed output type '%s'", reflect.TypeOf(content).Kind())
		}
	}
	return content
}

func (p *mapifier[T]) String(opts ...internal.OutputOpt) string {
	content := p.Out()

	return internal.NewOutput(opts...).String(content)
}

func (p *mapifier[T]) YAML() string {
	return internal.NewOutput().YAML(p.Out())
}

func (p *mapifier[T]) JSON() string {
	return internal.NewOutput().JSON(p.Out())
}

type setter[T Type] struct {
	s        *internal.Setter
	mapifier *mapifier[T]
}

func (s *setter[T]) Set(args ...any) Mapifier[T] {
	pathValueList := s.mapifier.sanitizer.SanitizePathValueList(args...)
	s.mapifier.mutators = append(s.mapifier.mutators, internal.NewSetter().Set(s.mapifier.parser, pathValueList)...)
	return s.mapifier
}
