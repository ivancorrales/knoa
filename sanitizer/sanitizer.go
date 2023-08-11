package sanitizer

import (
	"log"
)

var emptyValue = struct{}{}

type PathValue struct {
	Path  string
	Value any
}

type PathValueList []PathValue

func SanitizePathValueList(strict bool, args ...any) PathValueList {
	if len(args)%2 != 0 {
		args = append(args, emptyValue)
	}
	//nolint: gomnd
	list := make(PathValueList, len(args)/2)
	arg := 0
	invalidPathValues := 0
	for i := 0; i < len(args); i += 2 {
		path, ok := args[i].(string)
		if !ok {
			if strict {
				log.Panicf("invalid Path '%v'.  Paths must be string", args[i])
			}
			invalidPathValues += 1
			continue
		}
		list[arg] = PathValue{
			Path:  path,
			Value: args[i+1],
		}
		arg++
	}
	if invalidPathValues > 0 {
		return list[:len(list)-invalidPathValues]
	}
	return list
}

func SanitizePathList(strict bool, args ...any) []string {
	list := make([]string, len(args))
	invalidPathValues := 0
	for i := 0; i < len(args); i++ {
		path, ok := args[i].(string)
		if !ok {
			if strict {
				log.Panicf("invalid Path '%v'.  Paths must be string", args[i])
			}
			invalidPathValues += 1
			continue
		}
		list[i] = path
	}
	if invalidPathValues > 0 {
		return list[:len(list)-invalidPathValues]
	}
	return list
}
