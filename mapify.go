package mapify

import (
	"github.com/ivancorrales/mapify/mapifier"
)

func LoadMap(content map[string]any, opts ...mapifier.MapifyOpt) mapifier.Mapifier[map[string]any] {
	return mapifier.Load[map[string]any](content, opts...)
}

func Map(opts ...mapifier.MapifyOpt) mapifier.Mapifier[map[string]any] {
	return mapifier.New[map[string]any](opts...)
}

func LoadArray(content []any, opts ...mapifier.MapifyOpt) mapifier.Mapifier[[]any] {
	return mapifier.Load[[]any](content, opts...)
}

func Array(opts ...mapifier.MapifyOpt) mapifier.Mapifier[[]any] {
	return mapifier.New[[]any](opts...)
}
