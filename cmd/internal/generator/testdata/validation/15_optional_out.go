// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v T15Validator) Validate() error {
	if v.F1 != "" && !valid.Email(v.F1) {
		return errors.New("F1 must be a valid email address")
	}
	if v.F2 != "" {
		if !valid.Email(v.F2) {
			return errors.New("F2 must be a valid email address")
		} else if len(v.F2) < 5 || len(v.F2) > 85 {
			return errors.New("F2 must be of length between: 5 and 85 (inclusive)")
		}
	}
	if v.F3 != nil && *v.F3 != "" && !valid.Email(*v.F3) {
		return errors.New("F3 must be a valid email address")
	}
	if v.F4 != nil && *v.F4 != "" {
		if !valid.Email(*v.F4) {
			return errors.New("F4 must be a valid email address")
		} else if len(*v.F4) < 5 || len(*v.F4) > 85 {
			return errors.New("F4 must be of length between: 5 and 85 (inclusive)")
		}
	}
	return nil
}