// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"strings"
)

func (v Validator) Validate() error {
	v.F1 = strings.Replace(v.F1, "a", "b", 3)
	if v.F2 != nil {
		*v.F2 = strings.Replace(*v.F2, "foo", "bar", -1)
	}
	return nil
}
