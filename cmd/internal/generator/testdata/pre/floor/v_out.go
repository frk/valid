// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"math"
)

func (v Validator) Validate() error {
	v.F1 = math.Floor(v.F1)
	if v.F2 != nil {
		*v.F2 = math.Floor(*v.F2)
	}
	return nil
}
