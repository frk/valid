// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v T31Validator) Validate() error {
	if !valid.MAC(v.F1, 0) {
		return errors.New("F1 must be a valid MAC")
	}
	if v.F2 != nil && *v.F2 != nil && !valid.MAC(**v.F2, 6) {
		return errors.New("F2 must be a valid MAC")
	}
	if v.F3 == nil || *v.F3 == nil || **v.F3 == "" {
		return errors.New("F3 is required")
	} else if !valid.MAC(**v.F3, 8) {
		return errors.New("F3 must be a valid MAC")
	}
	return nil
}