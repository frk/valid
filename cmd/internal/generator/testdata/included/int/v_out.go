// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v Validator) Validate() error {
	if !valid.Int(v.F1) {
		return errors.New("F1 string content must match an integer")
	}
	if v.F2 != nil && !valid.Int(*v.F2) {
		return errors.New("F2 string content must match an integer")
	}
	return nil
}