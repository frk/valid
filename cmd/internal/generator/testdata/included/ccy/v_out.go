// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"

	"github.com/frk/valid"
)

func (v Validator) Validate() error {
	if !valid.Currency(v.F1, "usd", nil) {
		return errors.New("F1 must be a valid currency amount")
	}
	if v.F2 != nil && !valid.Currency(*v.F2, "usd", nil) {
		return errors.New("F2 must be a valid currency amount")
	}
	if !valid.Currency(v.F3, "gbp", nil) {
		return errors.New("F3 must be a valid currency amount")
	}
	if v.F4 != nil && !valid.Currency(*v.F4, "eur", v.ccyOpts) {
		return errors.New("F4 must be a valid currency amount")
	}
	if !valid.Currency(v.F5, "usd", v.ccyOpts) {
		return errors.New("F5 must be a valid currency amount")
	}
	return nil
}