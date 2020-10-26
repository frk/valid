// DO NOT EDIT. This file was generated by "github.com/frk/isvalid".

package testdata

import (
	"errors"

	"github.com/frk/isvalid"
)

func (v IPValidator) Validate() error {
	if !isvalid.IP(v.F1) {
		return errors.New("F1 must be a valid IP")
	}
	if v.F2 != nil && *v.F2 != nil {
		f := **v.F2
		if !isvalid.IP(f, 4) {
			return errors.New("F2 must be a valid IP")
		}
	}
	if v.F3 == nil || *v.F3 == nil || len(**v.F3) == 0 {
		return errors.New("F3 is required")
	} else {
		f := **v.F3
		if !isvalid.IP(f, 6, 4) {
			return errors.New("F3 must be a valid IP")
		}
	}
	return nil
}