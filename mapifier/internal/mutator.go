package internal

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/fatih/structs"
)

type MutatorOperation int32

const (
	SetOp MutatorOperation = iota
	UnsetOp
)

type Mutator struct {
	name      string
	index     string
	child     *Mutator
	value     any
	operation MutatorOperation
}

func (m *Mutator) IsArray() bool {
	_, err := strconv.Atoi(m.index)
	return err == nil || m.index == "*"
}

func (m *Mutator) Child() *Mutator {
	m.child.operation = m.operation
	return m.child
}

func (m *Mutator) Operation() MutatorOperation {
	return m.operation
}

func (m *Mutator) applyValue(in any) any {
	val := reflect.ValueOf(m.value)
	switch val.Kind() {
	case reflect.Struct:
		return structs.Map(val.Interface())
	case reflect.Func:
		x := reflect.TypeOf(m.value)
		if x.NumIn() != 1 || x.NumOut() != 1 {
			return nil
		}
		inVal := reflect.ValueOf(in)
		if x.In(0).Kind() != inVal.Kind() {
			return nil
		}
		out := val.Call([]reflect.Value{inVal})
		return out[0].String()
	case reflect.Slice, reflect.Array:
		out := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			if val.Index(i).Kind() == reflect.Struct {
				out[i] = structs.Map(val.Index(i).Interface())
			} else {
				out[i] = val.Index(i).Interface()
			}
		}
		return out
	default:
		return m.value
	}
}

func (m *Mutator) String() string {
	return m.pretty("")
}

func (m *Mutator) pretty(prefix string) string {
	// out := fmt.Sprintf("name: %s index: %s Kind: %v Value: %v", m.name, m.index, m.kind, m.value)
	out := fmt.Sprintf("name: %s index: %s  Value: %v", m.name, m.index, m.value)
	if m.child != nil {
		return fmt.Sprintf("%s \n%s %s ", out, prefix, m.child.pretty(prefix+"\t"))
	}
	return out
}

func (m *Mutator) WithValue(value any) {
	if m.child == nil {
		m.value = value
	} else {
		m.Child().WithValue(value)
	}
}

func (m *Mutator) addToBottom(child *Mutator) {
	if m.child == nil {
		m.child = child
	} else {
		m.Child().addToBottom(child)
	}
}

func (m *Mutator) ToMap(content map[string]any) map[string]any {
	if content == nil {
		content = make(map[string]any)
	}
	if m.child == nil {
		if m.operation == UnsetOp {
			delete(content, m.name)
			return content
		}
		if m.value != nil {
			content[m.name] = m.applyValue(content[m.name])
		}
		return content
	}
	mt := *m.Child()
	c := content[m.name]
	switch reflect.ValueOf(c).Kind() {
	case reflect.Slice, reflect.Array:
		var childContent []any
		if c == nil {
			childContent = make([]any, 0)
		} else {
			childContent, _ = c.([]any)
		}
		content[m.name] = mt.ToArray(childContent)
	default:
		var childContent map[string]any
		if c == nil {
			childContent = make(map[string]any, 0)
		} else {
			childContent, _ = c.(map[string]any)
		}
		mt.ToMap(childContent)
		if m.operation == UnsetOp {
			delete(content, m.name)
			return content
		}
		content[m.name] = childContent
	}

	return content
}

func (m *Mutator) ToArray(content []any) []any {
	if content == nil {
		content = make([]any, 0)
	}
	content = ensureSizeOfArray(content, m.index)

	index, err := strconv.Atoi(m.index)
	if err != nil {
		if m.index == "*" {
			for i := 0; i < len(content); i++ {
				content = m.itemToArray(i, content)
			}
		}
		return content
	}
	return m.itemToArray(index, content)
}

//nolint:nestif
func (m *Mutator) itemToArray(index int, content []any) []any {
	if m.child == nil {
		if m.operation == UnsetOp {
			return append(content[:index], content[index+1:]...)
		}
		if m.value != nil && index < len(content) {
			content[index] = m.applyValue(content[index])
		}
		return content
	}
	currentChild := content[index]
	if m.child.IsArray() {
		var childContent []any
		if currentChild == nil {
			childContent = make([]any, 0)
		} else {
			childContent, _ = currentChild.([]any)
		}
		content[index] = m.Child().ToArray(childContent)
	} else {
		var childContent map[string]any
		if currentChild == nil {
			childContent = make(map[string]any)
		} else {
			childContent, _ = currentChild.(map[string]any)
		}
		content[index] = m.Child().ToMap(childContent)
	}
	return content
}

func ensureSizeOfArray(arrayContent []any, indexStr string) []any {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return arrayContent
	}
	increment := index - len(arrayContent) + 1
	if increment > 0 {
		appendedArray := make([]any, increment)
		arrayContent = append(arrayContent, appendedArray...)
	}
	return arrayContent
}
