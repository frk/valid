// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"strings"
)

func (v Validator) Validate() error {
	v.F1 = strings.ToTitle(v.F1)
	if v.F2 != nil {
		*v.F2 = strings.ToTitle(*v.F2)
	}
	return nil
}
