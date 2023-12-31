package internal

import (
	"reflect"

	"github.com/fatih/structs"
)

func Normalize[T Type](input T) any {
	return normalize(input)
}

func normalize[T Type](input T) any {
	value := reflect.ValueOf(input)
	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		itemsLen := value.Len()
		if itemsLen == 0 {
			return make([]any, 1)
		}
		output := make([]any, itemsLen)
		for i := 0; i < itemsLen; i++ {
			itemValue := reflect.ValueOf(input).Index(i).Interface()
			output[i] = evalValue(itemValue)
		}
		return output
	case reflect.Struct:
		return structs.Map(input)
	case reflect.Map:
		output := make(map[string]any)
		if in, ok := reflect.ValueOf(input).Interface().(map[string]any); ok {
			for k, v := range in {
				output[k] = evalValue(v)
			}
		}
		return output
	}
	return input
}

func evalValue(in any) (out any) {
	switch reflect.ValueOf(in).Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		out = normalize(in)
	case reflect.Struct:
		out = normalize(structs.Map(in))
	default:
		out = in
	}
	return out
}
