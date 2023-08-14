package sanitizer

import (
	"log"
	"reflect"
)

var emptyValue = struct{}{}

var emptyFunc = func(in any)any{
	return in
}

type PathValue struct {
	Path  string
	Value any
}

type PathFunc struct {
	Path  string
	Func reflect.Value
}

type PathValueList []PathValue

type PathFuncList []PathFunc

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

func SanitizePathFuncList(strict bool, args ...any) PathFuncList {
	if len(args)%2 != 0 {
		args = append(args, emptyFunc)
	}
	//nolint: gomnd
	list := make(PathFuncList, len(args)/2)
	arg := 0
	invalidPathFuncs := 0
	for i := 0; i < len(args); i += 2 {
		path, ok := args[i].(string)
		if !ok {
			if strict {
				log.Panicf("invalid Path '%v'.  Paths must be string", args[i])
			}
			invalidPathFuncs += 1
			continue
		}
		fn:=reflect.ValueOf(args[i+1])
		if (fn.Kind()!=reflect.Func){
			if strict {
				log.Panicf("invalid Func '%v'.  Paths must be a valid func(any)any", args[i+1])
			}
			invalidPathFuncs += 1
			continue
		}
		
		list[arg] = PathFunc{
			Path:  path,
			Func: fn,
		}
		arg++
	}
	if invalidPathFuncs > 0 {
		return list[:len(list)-invalidPathFuncs]
	}
	return list
}