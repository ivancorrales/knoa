package knoa

import (
	"reflect"

	"github.com/ivancorrales/knoa/mapifier"
)

func LoadMap(content map[string]any, opts ...mapifier.KnoaOpt) mapifier.Mapifier[map[string]any] {
	return mapifier.Load[map[string]any](content, opts...)
}

func Map(opts ...mapifier.KnoaOpt) mapifier.Mapifier[map[string]any] {
	return mapifier.New[map[string]any](opts...)
}

func LoadArray(content []any, opts ...mapifier.KnoaOpt) mapifier.Mapifier[[]any] {
	return mapifier.Load[[]any](content, opts...)
}

func Array(opts ...mapifier.KnoaOpt) mapifier.Mapifier[[]any] {
	return mapifier.New[[]any](opts...)
}

func Load[T mapifier.Type](content T, opts ...mapifier.KnoaOpt) mapifier.Mapifier[T] {
	v := reflect.ValueOf(content)
	i := v.Interface()
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return mapifier.Load[[]any](i.([]any), opts...).(mapifier.Mapifier[T])
	case reflect.Map:
		return mapifier.Load[map[string]any](i.(map[string]any), opts...).(mapifier.Mapifier[T])
	}
	return mapifier.New[map[string]any]().(mapifier.Mapifier[T])
}