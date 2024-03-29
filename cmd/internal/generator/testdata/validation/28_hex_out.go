// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v T28Validator) Validate() error {
	if !valid.Hex(v.F1) {
		return errors.New("F1 must be a valid hexadecimal string")
	}
	if v.F2 != nil && *v.F2 != nil && !valid.Hex(**v.F2) {
		return errors.New("F2 must be a valid hexadecimal string")
	}
	if v.F3 == nil || *v.F3 == nil || **v.F3 == "" {
		return errors.New("F3 is required")
	} else if !valid.Hex(**v.F3) {
		return errors.New("F3 must be a valid hexadecimal string")
	}
	return nil
}
