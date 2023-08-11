package mutator

import (
	"github.com/fatih/structs"
	"reflect"
	"strconv"
)

type Mutator struct {
	name      string
	index     string
	child     *Mutator
	value     any
	operation operationCode
}

func (m *Mutator) addValueToNode(v any) {
	if m.child == nil {
		m.value = v
	} else {
		m.child.addValueToNode(v)
	}
}

func (m *Mutator) Child() *Mutator {
	m.child.operation = m.operation
	return m.child
}

func (m *Mutator) IsArray() bool {
	_, err := strconv.Atoi(m.index)
	return err == nil || m.index == "*"
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

func (m *Mutator) ToMap(content map[string]any) (map[string]any, error) {
	if content == nil {
		content = make(map[string]any)
	}
	if m.child == nil {
		if m.operation == unsetOp {
			delete(content, m.name)
			return content, nil
		}
		if m.value != nil {
			content[m.name] = m.applyValue(content[m.name])
		}
		return content, nil
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
		value, toArrayErr := mt.ToArray(childContent)
		if toArrayErr != nil {
			return nil, toArrayErr
		}
		content[m.name] = value
	default:
		var childContent map[string]any
		if c == nil {
			childContent = make(map[string]any)
		} else {
			childContent, _ = c.(map[string]any)
		}
		mt.ToMap(childContent)
		if m.operation == unsetOp {
			delete(content, m.name)
			return content, nil
		}
		content[m.name] = childContent
	}

	return content, nil
}

func (m *Mutator) ToArray(content []any) ([]any, error) {
	if content == nil {
		content = make([]any, 0)
	}
	content = ensureSizeOfArray(content, m.index)

	index, err := strconv.Atoi(m.index)
	if err != nil {
		if m.index == "*" {
			var itemToArrayErr error
			for i := 0; i < len(content); i++ {
				content, itemToArrayErr = m.itemToArray(i, content)
				if itemToArrayErr != nil {
					return nil, itemToArrayErr
				}
			}
			return content, nil
		}
		return nil, err
	}
	return m.itemToArray(index, content)
}

//nolint:nestif
func (m *Mutator) itemToArray(index int, content []any) ([]any, error) {
	if m.child == nil {
		if m.operation == unsetOp {
			return append(content[:index], content[index+1:]...), nil
		}
		if m.value != nil && index < len(content) {
			content[index] = m.applyValue(content[index])
		}
		return content, nil
	}
	if m.child.IsArray() {
		child := castOrCreateArray(content[index])
		c, err := m.Child().ToArray(child)
		if err != nil {
			return nil, err
		}
		content[index] = c
	} else {
		child := castOrCreateMap(content[index])
		mapValue, mapErr := m.Child().ToMap(child)
		if mapErr != nil {
			return nil, mapErr
		}
		content[index] = mapValue
	}
	return content, nil
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

func castOrCreateArray(in any) (out []any) {
	if in != nil {
		out, _ = in.([]any)
	}
	return

}
func castOrCreateMap(in any) (out map[string]any) {
	if in != nil {
		out, _ = in.(map[string]any)
	}
	return

}
