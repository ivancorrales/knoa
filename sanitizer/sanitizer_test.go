package sanitizer

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sanitizer_sanitize(t *testing.T) {
	type fields struct {
		strict bool
	}
	type args struct {
		args []any
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		panicked bool
		want     PathValueList
	}{
		{
			name: "The list of args are correctly provide, there's nothing to be sanitized",
			fields: fields{
				strict: true,
			},
			args: args{
				args: []any{"key1", true, "parent.key2", 20},
			},
			want: PathValueList{
				{"key1", true},
				{"parent.key2", 20},
			},
		},
		{
			name: "The list of args is empty",
			fields: fields{
				strict: true,
			},
			args: args{
				args: []any{},
			},
			want: PathValueList{},
		},
		{
			name: "The list of args is nil",
			fields: fields{
				strict: true,
			},
			args: args{
				args: nil,
			},
			want: PathValueList{},
		},
		{
			name: "The list of args contains an even number of items, and Strict mode is enabled",
			fields: fields{
				strict: true,
			},
			args: args{
				args: []any{"key1", "home", "key2"},
			},
			want: PathValueList{
				{"key1", "home"},
				{"key2", emptyValue},
			},
		},
		{
			name: "The list of args contains an even number of items, and Strict mode is disabled",
			fields: fields{
				strict: false,
			},
			args: args{
				args: []any{"key1", "home", "key2"},
			},
			want: PathValueList{
				{"key1", "home"},
				{"key2", emptyValue},
			},
		},
		{
			name: "The list of args contains a non string key, and enabled mode is disabled",
			fields: fields{
				strict: false,
			},
			args: args{
				args: []any{20, "home", "key2"},
			},
			want: PathValueList{
				{"key2", emptyValue},
			},
		},
		{
			name: "The list of args contains multiple  non string keys, and enabled mode is disabled",
			fields: fields{
				strict: false,
			},
			args: args{
				args: []any{20, "home", "key2", "hello", true, 20},
			},
			want: PathValueList{
				{"key2", "hello"},
			},
		},
		{
			name: "The list of args contains non string keys, and enabled mode is enabled",
			fields: fields{
				strict: true,
			},
			args: args{
				args: []any{20, "home", "key2"},
			},
			panicked: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.panicked {
				assert.Panics(t, func() { SanitizePathValueList(tt.fields.strict, tt.args.args...) }, "The execution should end panicking")
			} else {
				if got := SanitizePathValueList(tt.fields.strict, tt.args.args...); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SanitizePathValueList() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
