// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"
)

func (v T17Validator) Validate() error {
	if v.F1 == "" {
		return errors.New("F1 is required")
	}
	return nil
}
