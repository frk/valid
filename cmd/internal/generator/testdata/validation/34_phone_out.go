// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v T34Validator) Validate() error {
	if !valid.Phone(v.F1, "us") {
		return errors.New("F1 must be a valid phone number")
	}
	if v.F2 != nil && *v.F2 != nil && !valid.Phone(**v.F2, "us") {
		return errors.New("F2 must be a valid phone number")
	}
	if v.F3 == nil || *v.F3 == nil || **v.F3 == "" {
		return errors.New("F3 is required")
	} else if !valid.Phone(**v.F3, "jp") {
		return errors.New("F3 must be a valid phone number")
	}
	return nil
}