// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v Validator) Validate() error {
	if !valid.HSL(v.F1) {
		return errors.New("F1 must be a valid HSL color")
	}
	if v.F2 != nil && !valid.HSL(*v.F2) {
		return errors.New("F2 must be a valid HSL color")
	}
	return nil
}
