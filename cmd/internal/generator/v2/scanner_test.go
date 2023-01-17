package generate

import (
	"fmt"
	"testing"

	"github.com/frk/compare"
)

func Test_scanner(t *testing.T) {
	tests := []struct {
		in   string
		want []*token
	}{{
		in: "${0:+}!=$1",
		want: []*token{
			{t: t_paramx, p: 0, v: "0:+"},
			{t: t_neq, p: 6},
			{t: t_param, p: 8, v: "1"},
			{t: t_eof}},
	}, {
		in: "  ${0:+} !=	 $1\t",
		want: []*token{
			{t: t_ws, p: 0, v: "  "},
			{t: t_paramx, p: 2, v: "0:+"},
			{t: t_ws, p: 8, v: " "},
			{t: t_neq, p: 9, v: ""},
			{t: t_ws, p: 11, v: "	 "},
			{t: t_param, p: 13, v: "1"},
			{t: t_ws, p: 15, v: "\t"},
			{t: t_eof}},
	}, {
		in: "foo",
		want: []*token{
			{t: t_ident, v: "foo"},
			{t: t_eof}},
	}, {
		in: "if ${0} > ${1} {",
		want: []*token{
			{t: t_if, p: 0, v: ""},
			{t: t_ws, p: 2, v: " "},
			{t: t_paramx, p: 3, v: "0"},
			{t: t_ws, p: 7, v: " "},
			{t: t_gtr, p: 8, v: ""},
			{t: t_ws, p: 9, v: " "},
			{t: t_paramx, p: 10, v: "1"},
			{t: t_ws, p: 14, v: " "},
			{t: t_lbrace, p: 15, v: ""},
			{t: t_eof}},
	}}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			out := make(chan *token)
			go (&scanner{in: tt.in, out: out}).run()

			got := make([]*token, 0)
			for t := range out {
				if false { // debug?
					fmt.Printf("%#v\n", t)
				}

				got = append(got, t)
				if t.t == t_eof || t.t == t_invalid {
					break
				}
			}
			if e := compare.Compare(got, tt.want); e != nil {
				t.Errorf("Error: %v", e)
			}
		})
	}
}
