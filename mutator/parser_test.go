package mutator

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parser_parse(t *testing.T) {
	type fields struct {
		strict bool
	}
	type args struct {
		pathExpr string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     *Mutator
		panicked bool
	}{
		{
			name: "A single attribute name",
			fields: fields{
				strict: false,
			},
			args: args{
				pathExpr: "fullnameOverride",
			},
			want: &Mutator{
				child: &Mutator{
					name: "fullnameOverride",
				},
			},
		},
		{
			name: "A single Array",
			fields: fields{
				strict: false,
			},
			args: args{
				pathExpr: "[1]",
			},
			want: &Mutator{
				child: &Mutator{
					index: "1",
				},
			},
		},
		{
			name: "two deep level valid expression Path",
			fields: fields{
				strict: false,
			},
			args: args{
				pathExpr: "people[0].firstname",
			},
			want: &Mutator{
				child: &Mutator{
					name: "people",

					child: &Mutator{
						index: "0",

						child: &Mutator{
							name: "firstname",
						},
					},
				},
			},
		},
		{
			name: "An invalid expression but Strict mode is disabled",
			fields: fields{
				strict: false,
			},
			args: args{
				pathExpr: "peopl\\\\e[0].firstname",
			},
			want: nil,
		},
		{
			name: "An invalid expression and Strict mode is enabled",
			fields: fields{
				strict: true,
			},
			args: args{
				pathExpr: "peopl\\\\e[0].firstname",
			},
			want:     nil,
			panicked: true,
		},
		{
			name: "A simple Array",
			fields: fields{
				strict: true,
			},
			args: args{
				pathExpr: "[0]",
			},
			want: &Mutator{
				child: &Mutator{
					index: "0",
				},
			},
		},
		{
			name: "Multiple arrays",
			fields: fields{
				strict: true,
			},
			args: args{
				pathExpr: "[0][1][2].name",
			},
			want: &Mutator{
				child: &Mutator{
					index: "0",
					child: &Mutator{
						index: "1",
						child: &Mutator{
							index: "2",

							child: &Mutator{
								name: "name",
							},
						},
					},
				},
			},
		},
		{
			name: "single Array",
			fields: fields{
				strict: true,
			},
			args: args{
				pathExpr: "[2]",
			},
			want: &Mutator{
				child: &Mutator{
					index: "2",
				},
			},
		},
		{
			name: "Attributes contains dots ",
			fields: fields{
				strict: false,
			},
			args: args{
				pathExpr: "annotations.\"a.b.c\"",
			},
			want: &Mutator{
				child: &Mutator{
					name: "annotations",

					child: &Mutator{
						name: "a.b.c",
					},
				},
			},
		},
		{
			name: "Attributes in the middle of a Path contains dots ",
			fields: fields{
				strict: false,
			},
			args: args{
				pathExpr: "annotations.\"a.b.c\".Value[0].name",
			},
			want: &Mutator{
				child: &Mutator{
					name: "annotations",

					child: &Mutator{
						name: "a.b.c",

						child: &Mutator{
							name: "Value",

							child: &Mutator{
								index: "0",

								child: &Mutator{
									name: "name",
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pathRegExp, attrRegExp := RegExpsFromAttributeFormat(DefAttributeNameFormat)
			p := &Parser{
				RegExp:          pathRegExp,
				Strict:          tt.fields.strict,
				AttributeRegExp: attrRegExp,
			}

			if tt.panicked {
				assert.Panics(t, func() { p.Parse(tt.args.pathExpr) }, "The execution should end panicking")
			} else {
				res, _ := p.Parse(tt.args.pathExpr)
				assertParsedElements(t, tt.want, res)
			}
		})
	}
}

func assertParsedElements(t *testing.T, expected *Mutator, got *Mutator) {
	if expected == nil && got == nil {
		return
	}
	if (expected == nil && got != nil) || (expected != nil && got == nil) {
		t.Errorf("\nexpected= %v  , \ngot= %v", expected, got)
		return
	}
	if (expected.child == nil && got.child != nil) || (expected.child != nil && got.child == nil) {
		t.Errorf("\nexpected= %v  , \ngot= %v", expected, got)
		return
	}
	if got.child != nil {
		assertParsedElements(t, expected.child, got.child)
		return
	}
	if expected.name != got.name || expected.value != got.value || expected.index != got.index {
		t.Errorf("\nexpected= %v  , got= %v", expected, got)
		return
	}
}

func Test_RegExpFromAttributeFormat(t *testing.T) {
	type args struct {
		attributeFormat string
	}
	tests := []struct {
		name string
		args args
		want *regexp.Regexp
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, RegExpFromAttributeFormat(tt.args.attributeFormat), "RegExpFromAttributeFormat(%v)", tt.args.attributeFormat)
		})
	}
}
