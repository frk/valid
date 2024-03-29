// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v Validator) Validate() error {
	if !valid.ISSN(v.F1, true, true) {
		return errors.New("F1 must be a valid ISSN")
	}
	if v.F2 != nil && !valid.ISSN(*v.F2, false, false) {
		return errors.New("F2 must be a valid ISSN")
	}
	return nil
}
