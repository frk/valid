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
			{t: t_paramx, v: "0:+"},
			{t: t_neq},
			{t: t_param, v: "1"},
			{t: t_eof}},
	}, {
		in: "  ${0:+} !=	 $1\t",
		want: []*token{
			{t: t_ws},
			{t: t_paramx, v: "0:+"},
			{t: t_ws},
			{t: t_neq},
			{t: t_ws},
			{t: t_param, v: "1"},
			{t: t_ws},
			{t: t_eof}},
	}, {
		in: "foo",
		want: []*token{
			{t: t_ident, v: "foo"},
			{t: t_eof}},
	}, {
		in: "if ${0} > ${1} {",
		want: []*token{
			{t: t_if},
			{t: t_ws},
			{t: t_paramx, v: "0"},
			{t: t_ws},
			{t: t_gtr},
			{t: t_ws},
			{t: t_paramx, v: "1"},
			{t: t_ws},
			{t: t_lbrace},
			{t: t_eof}},
	}}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			c := make(chan *token)
			ok := make(chan bool)
			got := make([]*token, 0)
			go func() {
				for t := range c {
					got = append(got, t)
					if false { // debug?
						fmt.Printf("%#v\n", t)
					}
				}
				ok <- true
			}()

			(&scanner{in: tt.in, out: c}).run()
			close(c)

			if <-ok {
				if e := compare.Compare(got, tt.want); e != nil {
					t.Errorf("Error: %v", e)
				}
			}
		})
	}
}
