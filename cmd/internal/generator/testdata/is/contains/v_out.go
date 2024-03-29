// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"
	"strings"
)

func (v Validator) Validate() error {
	if !strings.Contains(v.F1, "foo") {
		return errors.New("F1 must contain substring: \"foo\"")
	}
	if v.F2 != nil && !strings.Contains(*v.F2, "bar") {
		return errors.New("F2 must contain substring: \"bar\"")
	}
	if !strings.Contains(v.F3, "foo") && !strings.Contains(v.F3, "bar") && !strings.Contains(v.F3, "baz") {
		return errors.New("F3 must contain substring: \"foo\" or \"bar\" or \"baz\"")
	}
	return nil
}
