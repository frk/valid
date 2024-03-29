// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"
	"strings"

	"github.com/frk/valid"
)

func (v T47Validator) Validate() error {
	if v.F7a != nil && !valid.Email(*v.F7a) {
		return errors.New("F7a must be a valid email address")
	}
	if v.F7b != nil && *v.F7b != nil && !valid.Email(**v.F7b) {
		return errors.New("F7b must be a valid email address")
	}
	if v.F8a != nil {
		if !valid.Hex(*v.F8a) {
			return errors.New("F8a must be a valid hexadecimal string")
		} else if len(*v.F8a) < 8 || len(*v.F8a) > 128 {
			return errors.New("F8a must be of length between: 8 and 128 (inclusive)")
		}
	}
	if v.F8b != nil && *v.F8b != nil {
		f := **v.F8b
		if !valid.Hex(f) {
			return errors.New("F8b must be a valid hexadecimal string")
		} else if len(f) < 8 || len(f) > 128 {
			return errors.New("F8b must be of length between: 8 and 128 (inclusive)")
		}
	}
	if v.F10 != nil && *v.F10 != nil && **v.F10 != nil && ***v.F10 != nil && ****v.F10 != nil {
		f := *****v.F10
		if !strings.HasPrefix(f, "foo") {
			return errors.New("F10 must be prefixed with: \"foo\"")
		} else if !strings.Contains(f, "bar") {
			return errors.New("F10 must contain substring: \"bar\"")
		} else if !strings.HasSuffix(f, "baz") && !strings.HasSuffix(f, "quux") {
			return errors.New("F10 must be suffixed with: \"baz\" or \"quux\"")
		} else if len(f) < 8 || len(f) > 64 {
			return errors.New("F10 must be of length between: 8 and 64 (inclusive)")
		}
	}
	return nil
}
