// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package testdata

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/frk/valid"
)

func (v T04Validator) Validate() error {
	v.F1 = strings.TrimSpace(v.F1)
	if v.F1 == "" {
		return errors.New("F1 is required")
	} else if !valid.Email(v.F1) {
		return errors.New("F1 must be a valid email address")
	}
	v.F2 = strings.Repeat(v.F2, 2)
	if v.F2 == "" {
		return errors.New("F2 is required")
	} else if len(v.F2) > 256 {
		return errors.New("F2 must be of length at most: 256")
	}
	if v.F3 == nil {
		return errors.New("F3 is required")
	} else {
		*v.F3 = strconv.Quote(strings.ToTitle(strings.TrimSpace(*v.F3)))
		if *v.F3 == "" {
			return errors.New("F3 is required")
		} else if !valid.Alnum(*v.F3, "en") {
			return errors.New("F3 must be an alphanumeric string")
		}
	}
	if v.F4 == nil || *v.F4 == nil || **v.F4 == nil {
		return errors.New("F4 is required")
	} else {
		***v.F4 = strings.ToUpper(***v.F4)
		f := ***v.F4
		if f == "" {
			return errors.New("F4 is required")
		} else if !valid.UpperCase(f) {
			return errors.New("F4 must contain only upper-case characters")
		}
	}
	if v.F5 == nil || *v.F5 == nil {
		return errors.New("F5 is required")
	} else {
		**v.F5 = math.Round(**v.F5)
		if **v.F5 == 0.0 {
			return errors.New("F5 is required")
		}
	}
	return nil
}