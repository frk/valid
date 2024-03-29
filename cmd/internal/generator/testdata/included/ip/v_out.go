// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v Validator) Validate() error {
	if !valid.IP(v.F1, 0) {
		return errors.New("F1 must be a valid IP")
	}
	if v.F2 != nil && !valid.IP(*v.F2, 0) {
		return errors.New("F2 must be a valid IP")
	}
	if !valid.IP(v.F3, 4) {
		return errors.New("F3 must be a valid IP")
	}
	if v.F4 != nil && !valid.IP(*v.F4, 6) {
		return errors.New("F4 must be a valid IP")
	}
	return nil
}
