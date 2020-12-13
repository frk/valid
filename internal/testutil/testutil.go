package testutil

import (
	"testing"
)

type List []struct {
	Name string
	Func func(v, cc string) bool
	Pass []string
	Fail []string
}

func Run(t *testing.T, ccs []string, list List) {
	for _, tt := range list {
		for _, cc := range ccs {
			name := cc + "/" + tt.Name
			for _, v := range tt.Pass {
				want := true
				t.Run(name, func(t *testing.T) {
					got := tt.Func(v, cc)
					if got != want {
						t.Errorf("got=%t; want=%t; value=%q", got, want, v)
					}
				})
			}
			for _, v := range tt.Fail {
				want := false
				t.Run(name, func(t *testing.T) {
					got := tt.Func(v, cc)
					if got != want {
						t.Errorf("got=%t; want=%t; value=%q", got, want, v)
					}
				})
			}
		}
	}
}
