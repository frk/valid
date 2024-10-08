// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"strings"
)

func (v Validator) Validate() error {
	v.F1 = strings.TrimSpace(v.F1)
	if v.F2 != nil {
		*v.F2 = strings.TrimSpace(*v.F2)
	}
	v.F3 = strings.TrimRight(strings.TrimLeft(v.F3, "abc"), "xyz")
	for _, e1 := range v.F4.A {
		if e1.X != nil {
			*e1.X = strings.TrimRight(strings.TrimLeft(strings.TrimSpace(*e1.X), "abc"), "xyz")
		}
	}
	return nil
}
