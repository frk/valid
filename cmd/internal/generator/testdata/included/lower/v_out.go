// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v Validator) Validate() error {
	if !valid.LowerCase(v.F1) {
		return errors.New("F1 must contain only lower-case characters")
	}
	if v.F2 != nil && !valid.LowerCase(*v.F2) {
		return errors.New("F2 must contain only lower-case characters")
	}
	return nil
}