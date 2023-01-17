package generate

import (
	"strings"
	"testing"
)

func Test_parser(t *testing.T) {
	t.Skip()

	type input struct {
		in   string
		args []any
	}

	tests := []struct {
		name  string
		input []input
		want  []string
		debug bool
	}{{
		name:  "empty block",
		input: []input{},
		want:  []string{},
	}, {
		name:  "lit_expr",
		input: []input{{in: "123"}},
		want:  []string{"123"},
	}, {
		name:  "lit_expr",
		input: []input{{in: "0.123"}},
		want:  []string{"0.123"},
	}, {
		name:  "lit_expr",
		input: []input{{in: "10_000"}},
		want:  []string{"10_000"},
	}, {
		name:  "ident_expr",
		input: []input{{in: "x"}},
		want:  []string{"x"},
	}, {
		name:  "ident_expr",
		input: []input{{in: "true"}},
		want:  []string{"true"},
	}, {
		name:  "ident_expr",
		input: []input{{in: "nil"}},
		want:  []string{"nil"},
	}, {
		name:  "unary_expr",
		input: []input{{in: "!x"}},
		want:  []string{"!x"},
	}, {
		name:  "unary_expr",
		input: []input{{in: "*x"}},
		want:  []string{"*x"},
	}, {
		name:  "unary_expr",
		input: []input{{in: "***x"}},
		want:  []string{"***x"},
	}, {
		name:  "unary_expr",
		input: []input{{in: "&x"}},
		want:  []string{"&x"},
	}, {
		name:  "binary_expr",
		input: []input{{in: "x > y"}},
		want:  []string{"x > y"},
	}, {
		name:  "binary_expr",
		input: []input{{in: "x && y"}},
		want:  []string{"x && y"},
	}, {
		name:  "binary_expr",
		input: []input{{in: "99 > 9.9"}},
		want:  []string{"99 > 9.9"},
	}, {
		name:  "binary_expr",
		input: []input{{in: "true || false"}},
		want:  []string{"true || false"},
	}, {
		name:  "index_expr",
		input: []input{{in: "x[i]"}},
		want:  []string{"x[i]"},
	}, {
		name:  "index_expr",
		input: []input{{in: "x[*p]"}},
		want:  []string{"x[*p]"},
	}, {
		name:  "call_expr",
		input: []input{{in: "x()"}},
		want:  []string{"x()"},
	}, {
		name:  "paren_expr",
		input: []input{{in: "()"}},
		want:  []string{"()"},
	}, {
		name:  "paren_expr",
		input: []input{{in: "(x)"}},
		want:  []string{"(x)"},
	}, {
		name:  "paren_expr",
		input: []input{{in: "(x > y)"}},
		want:  []string{"(x > y)"},
	}, {
		name:  "paren_expr",
		input: []input{{in: "(foo || bar)"}},
		want:  []string{"(foo || bar)"},
	}, {
		name:  "paren_expr",
		input: []input{{in: "(x) && (y)"}},
		want:  []string{"(x) && (y)"},
	}, {
		name:  "paren_expr",
		input: []input{{in: "((foo) || bar)"}},
		want:  []string{"((foo) || bar)"},
	}, {
		name: "paren_expr",
		input: []input{
			{in: "(x < y || x == z) && foo"},
		},
		want: []string{"(x < y || x == z) && foo"},
	}, {
		name: "paren_expr",
		input: []input{
			{in: "((x < y || x == z) && (y > b && z <= a) && sky_is_blue)"},
		},
		want: []string{"((x < y || x == z) && (y > b && z <= a) && sky_is_blue)"},
	}, {
		name: "paren_expr",
		input: []input{
			{in: "(((((x))) < (y)) && (((a)) > b))"},
		},
		want: []string{"(((((x))) < (y)) && (((a)) > b))"},
	}, {
		name: "paren_expr",
		input: []input{
			{in: "x < (y) && (a > b && (air_is_cold)) || z == true"},
		},
		want: []string{"x < (y) && (a > b && (air_is_cold)) || z == true"},
	}, {
		name: "if_stmt",
		input: []input{
			{in: "if x {"},
		},
		want: []string{"if x {}"},
	}, {
		name: "if_stmt",
		input: []input{
			{in: "if v := x; v {"},
		},
		want: []string{"if v := x; v {}"},
	}, {
		name: "if_stmt",
		input: []input{
			{in: "if v1, v2 := x1, x2; (v1 < v2) && sky_is_blue {"},
		},
		want: []string{"if v1, v2 := x1, x2; (v1 < v2) && sky_is_blue {}"},
	}, {
		name: "if_stmt",
		input: []input{
			{in: "if v1, v2, v3 := x1, (x2 > x3), (z == y); v1 && (v2 || v3) {"},
		},
		want: []string{"if v1, v2, v3 := x1, (x2 > x3), (z == y); v1 && (v2 || v3) {}"},
	}, {
		name: "for_stmt",
		input: []input{
			{in: "for {"},
		},
		want: []string{"for {}"},
	}, {
		name: "for_stmt",
		input: []input{
			{in: "for x != nil {"},
		},
		want: []string{"for x != nil {}"},
	}, {
		name: "for_stmt",
		input: []input{
			{in: "for i := 0; i < n; i++ {"},
		},
		want: []string{"for i := 0; i < n; i++ {}"},
	}, {
		name: "range_stmt",
		input: []input{
			{in: "for k, v := range m {"},
		},
		want: []string{"for k, v := range m {}"},
	}, {
		name: "range_stmt",
		input: []input{
			{in: "for range m {"},
		},
		want: []string{"for range m {}"},
	}, {
		name: "range_stmt",
		input: []input{
			{in: "for x := range m {"},
		},
		want: []string{"for x := range m {}"},
	}, {
		name: "range_stmt",
		input: []input{
			{in: "for _, v := range m {"},
		},
		want: []string{"for _, v := range m {}"},
	}, {
		name: "range_stmt",
		input: []input{
			{in: "for k, _ = range m {"},
		},
		want: []string{"for k, _ = range m {}"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := new(parser).init()
			p.debug = tt.debug
			for _, v := range tt.input {
				p.parse(v.in, v.args...)
			}
			p.close()

			want := strings.Join(tt.want, "\n")
			got, err := node_string(p.body)
			if err != nil {
				t.Error(err)
			} else if err == nil {
				got = strings.TrimSpace(got)
				got = strings.TrimPrefix(got, "{")
				got = strings.TrimSuffix(got, "}")
				got = strings.TrimSpace(got)

				if want != got {
					t.Errorf("\nwant: %q\ngot:  %q", want, got)
				}
			}
		})
	}
}
