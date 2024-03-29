// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v Validator) Validate() error {
	if !valid.SSN(v.F1) {
		return errors.New("F1 must be a valid SSN")
	}
	if v.F2 != nil && !valid.SSN(*v.F2) {
		return errors.New("F2 must be a valid SSN")
	}
	return nil
}
