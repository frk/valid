// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v Validator) Validate() error {
	if !valid.MD5(v.F1) {
		return errors.New("F1 must be a valid MD5 hash")
	}
	if v.F2 != nil && !valid.MD5(*v.F2) {
		return errors.New("F2 must be a valid MD5 hash")
	}
	return nil
}