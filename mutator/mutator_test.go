package mutator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ensureSizeOfArray(t *testing.T) {
	type args struct {
		arrayContent []any
		indexStr     string
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "The size of the Array mustn't be changed",
			args: args{
				arrayContent: []any{10, 20, 30},
				indexStr:     "2",
			},
			want: []any{10, 20, 30},
		},
		{
			name: "The size of the Array must be changed",
			args: args{
				arrayContent: []any{10, 20, 30},
				indexStr:     "7",
			},
			want: []any{10, 20, 30, nil, nil, nil, nil, nil},
		},
		{
			name: "The index is invalid",
			args: args{
				arrayContent: []any{10, 20, 30},
				indexStr:     "A",
			},
			want: []any{10, 20, 30},
		},
		{
			name: "The index is a non-positive Value",
			args: args{
				arrayContent: []any{10, 20, 30},
				indexStr:     "-A",
			},
			want: []any{10, 20, 30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ensureSizeOfArray(tt.args.arrayContent, tt.args.indexStr), "ensureSizeOfArray(%v, %v)", tt.args.arrayContent, tt.args.indexStr)
		})
	}
}

func Test_mutator_toArray(t *testing.T) {
	type fields struct {
		name  string
		index string
		child *Mutator

		value any
	}
	type args struct {
		content []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []any
	}{
		{
			name: "Add a new entry into the Array",
			fields: fields{
				index: "1",
				child: &Mutator{
					name:  "firstname",
					value: "Mary",
				},
			},
			args: args{
				content: []any{
					map[string]any{
						"firstname": "Jane",
					},
				},
			},
			want: []any{
				map[string]any{
					"firstname": "Jane",
				},
				map[string]any{
					"firstname": "Mary",
				},
			},
		},
		{
			name: "Modify the Value of an existing item in the Array",
			fields: fields{
				index: "1",
				child: &Mutator{
					name:  "firstname",
					value: "Mary",
				},
			},
			args: args{
				content: []any{
					map[string]any{
						"firstname": "Jane",
					},
				},
			},
			want: []any{
				map[string]any{
					"firstname": "Jane",
				},
				map[string]any{
					"firstname": "Mary",
				},
			},
		},
		{
			name: "The initial Array is empty",
			fields: fields{
				index: "1",
				child: &Mutator{
					name:  "firstname",
					value: "Mary",
				},
			},
			args: args{
				content: []any{},
			},
			want: []any{
				nil,
				map[string]any{
					"firstname": "Mary",
				},
			},
		},
		{
			name: "Add a new entry into the Array that contains sub arrays",
			fields: fields{
				index: "1",
				child: &Mutator{
					index: "3",
					value: "hello",
				},
			},
			args: args{
				content: []any{
					[]any{
						"my friend",
					},
				},
			},
			want: []any{
				[]any{
					"my friend",
				},
				[]any{
					nil, nil, nil, "hello",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mutator{
				name:  tt.fields.name,
				index: tt.fields.index,
				child: tt.fields.child,
				value: tt.fields.value,
			}
			c, _ := m.ToArray(tt.args.content)
			assert.Equalf(t, tt.want, c, "ToArray(%v)", tt.args.content)
		})
	}
}

func Test_mutator_toMap(t *testing.T) {
	type fields struct {
		parentExpr string
		name       string
		index      string
		path       string
		child      *Mutator
		value      any
	}
	type args struct {
		content map[string]any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]any
	}{
		{
			name: "Add a new entry into the Array",
			fields: fields{
				name: "item2",
				child: &Mutator{
					name:  "firstname",
					value: "Mary",
				},
			},
			args: args{
				content: map[string]any{
					"item1": map[string]any{
						"firstname": "Jane",
					},
				},
			},
			want: map[string]any{
				"item1": map[string]any{
					"firstname": "Jane",
				},
				"item2": map[string]any{
					"firstname": "Mary",
				},
			},
		},
		{
			name: "Modify the Value of an existing item in the map",
			fields: fields{
				name: "item1",
				child: &Mutator{
					name:  "firstname",
					value: "Mary",
				},
			},
			args: args{
				content: map[string]any{
					"item1": map[string]any{
						"firstname": "Jane",
					},
				},
			},
			want: map[string]any{
				"item1": map[string]any{
					"firstname": "Mary",
				},
			},
		},
		{
			name: "The initial Array is empty",
			fields: fields{
				name: "item",
				child: &Mutator{
					name:  "firstname",
					value: "Mary",
				},
			},
			args: args{
				content: map[string]any{},
			},
			want: map[string]any{
				"item": map[string]any{
					"firstname": "Mary",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mutator{
				name:  tt.fields.name,
				index: tt.fields.index,
				child: tt.fields.child,
				value: tt.fields.value,
			}
			c, _ := m.ToMap(tt.args.content)
			assert.Equalf(t, tt.want, c, "ToMap(%v)", tt.args.content)
		})
	}
}

func Test_mutator_withValue(t *testing.T) {
	type fields struct {
		parentExpr string
		name       string
		index      string
		path       string
		child      *Mutator
		value      any
	}
	type args struct {
		value any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "AAdd a Value to the Mutator which sdoesn't contain any Value yet",
			fields: fields{
				value: nil,
			},
			args: args{
				value: 20,
			},
		},
		{
			name: "AAdd a Value to the Mutator and overwrite its Value",
			fields: fields{
				value: 21,
			},
			args: args{
				value: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mutator{
				name:  tt.fields.name,
				index: tt.fields.index,
				child: tt.fields.child,

				value: tt.fields.value,
			}
			m.value = tt.args.value
			assert.Equal(t, tt.args.value, m.value)
		})
	}
}
