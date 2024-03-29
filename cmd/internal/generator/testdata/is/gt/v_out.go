// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"
)

func (v Validator) Validate() error {
	if v.F1 <= 3.14 {
		return errors.New("F1 must be greater than: 3.14")
	}
	if v.F2 != nil && *v.F2 <= 314 {
		return errors.New("F2 must be greater than: 314")
	}
	if v.F3 == nil || *v.F3 == 0 {
		return errors.New("F3 is required")
	} else if *v.F3 <= 31 {
		return errors.New("F3 must be greater than: 31")
	}
	return nil
}
