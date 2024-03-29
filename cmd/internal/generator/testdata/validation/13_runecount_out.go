// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"
	"unicode/utf8"
)

func (v T13Validator) Validate() error {
	if utf8.RuneCountInString(v.F1) != 10 {
		return errors.New("F1 must have rune count: 10")
	}
	if v.F2 != nil && (utf8.RuneCountInString(*v.F2) < 8 || utf8.RuneCountInString(*v.F2) > 256) {
		return errors.New("F2 must have rune count between: 8 and 256 (inclusive)")
	}
	if v.F3 == nil || *v.F3 == nil || len(**v.F3) == 0 {
		return errors.New("F3 is required")
	} else if utf8.RuneCount(**v.F3) < 1 || utf8.RuneCount(**v.F3) > 2 {
		return errors.New("F3 must have rune count between: 1 and 2 (inclusive)")
	}
	if utf8.RuneCount(v.F4) < 4 {
		return errors.New("F4 must have rune count at least: 4")
	}
	if utf8.RuneCountInString(v.F5) > 15 {
		return errors.New("F5 must have rune count at most: 15")
	}
	return nil
}
