// DO NOT EDIT. This file was generated by "github.com/frk/valid".

package bench

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/frk/valid"
)

func init() {
	valid.RegisterRegexp(`^[a-z]+\[[0-9]+\]$`)
	valid.RegisterRegexp(`^[a-z]+\[[0-9]+\]$`)
}

func (v FrkValidator) Validate() error {
	if v.F01 == "" {
		return errors.New("F01 is required")
	}
	if v.F02 == nil || *v.F02 == "" {
		return errors.New("F02 is required")
	}
	if v.F03 != nil && *v.F03 != 3.14 {
		return errors.New("F03 must be equal to: 3.14")
	}
	if v.F04 != nil && *v.F04 <= 123 {
		return errors.New("F04 must be greater than: 123")
	}
	if v.F05 > 42 {
		return errors.New("F05 must be less than or equal to: 42")
	}
	if v.F06 < 3.14 || v.F06 > 42 {
		return errors.New("F06 must be between: 3.14 and 42")
	}
	if v.F07 < 3.14 {
		return errors.New("F07 must be greater than or equal to: 3.14")
	} else if v.F07 > 42 {
		return errors.New("F07 must be less than or equal to: 42")
	}
	if len(v.F08) < 8 || len(v.F08) > 256 {
		return errors.New("F08 must be of length between: 8 and 256 (inclusive)")
	}
	if len(v.F09) != 12 {
		return errors.New("F09 must be of length: 12")
	}
	if len(v.F10) > 12 {
		return errors.New("F10 must be of length at most: 12")
	}
	if utf8.RuneCount(v.F11) < 4 {
		return errors.New("F11 must have rune count at least: 4")
	}
	if utf8.RuneCountInString(v.F12) > 15 {
		return errors.New("F12 must have rune count at most: 15")
	}
	if v.F13 == nil || *v.F13 == "" {
		return errors.New("F13 is required")
	} else if !strings.HasPrefix(*v.F13, "hello") {
		return errors.New("F13 must be prefixed with: \"hello\"")
	}
	if v.F14 != nil && !strings.Contains(*v.F14, "hello") {
		return errors.New("F14 must contain substring: \"hello\"")
	}
	if !valid.Alnum(v.F15, "en") {
		return errors.New("F15 must be an alphanumeric string")
	}
	if !valid.CIDR(v.F16) {
		return errors.New("F16 must be a valid CIDR notation")
	}
	if v.F17 != nil && !valid.FQDN(*v.F17) {
		return errors.New("F17 must be a valid FQDN")
	}
	if v.F18 != nil && !valid.IP(*v.F18, 0) {
		return errors.New("F18 must be a valid IP")
	}
	if !valid.Email(v.F19) {
		return errors.New("F19 must be a valid email address")
	}
	if !valid.Phone(v.F20, "us") {
		return errors.New("F20 must be a valid phone number")
	}
	if v.Sub1 != nil {
		if len(v.Sub1.F01) > 10 {
			return errors.New("Sub1.F01 must be of length at most: 10")
		} else {
			for _, e := range v.Sub1.F01 {
				if !valid.Email(e) {
					return errors.New("Sub1.F01 must be a valid email address")
				}
			}
		}
		if len(v.Sub1.F02) == 0 {
			return errors.New("Sub1.F02 is required")
		} else {
			for _, e := range v.Sub1.F02 {
				if !valid.Phone(e, "us") {
					return errors.New("Sub1.F02 must be a valid phone number")
				}
			}
		}
		if v.Sub1.F03 == nil || len(*v.Sub1.F03) == 0 {
			return errors.New("Sub1.F03 is required")
		} else {
			for k := range *v.Sub1.F03 {
				if !valid.Email(k) {
					return errors.New("Sub1.F03 must be a valid email address")
				}
			}
		}
		for k, e := range v.Sub1.F04 {
			if !valid.Email(k) {
				return errors.New("Sub1.F04 must be a valid email address")
			}
			if !valid.Phone(e, "us") {
				return errors.New("Sub1.F04 must be a valid phone number")
			}
		}
		if len(v.Sub1.F05) < 5 {
			return errors.New("Sub1.F05 must be of length at least: 5")
		} else {
			for k, e := range v.Sub1.F05 {
				if !valid.Phone(k, "us") {
					return errors.New("Sub1.F05 must be a valid phone number")
				}
				if e < 21 || e > 99 {
					return errors.New("Sub1.F05 must be between: 21 and 99")
				}
			}
		}
		if v.Sub1.F06 != nil && !valid.Match(*v.Sub1.F06, `^[a-z]+\[[0-9]+\]$`) {
			return errors.New("Sub1.F06 must match the regular expression: \"^[a-z]+\\\\[[0-9]+\\\\]$\"")
		}
		if !strings.HasSuffix(v.Sub1.F07, "goodbye") {
			return errors.New("Sub1.F07 must be suffixed with: \"goodbye\"")
		}
		if v.Sub1.F08 != nil && *v.Sub1.F08 != v.Arg1 {
			return fmt.Errorf("Sub1.F08 must be equal to: %v", v.Arg1)
		}
		if v.Sub1.F09 >= v.Arg2 {
			return fmt.Errorf("Sub1.F09 must be less than: %v", v.Arg2)
		}
		if v.Sub1.F10 <= v.Arg3 {
			return fmt.Errorf("Sub1.F10 must be greater than: %v", v.Arg3)
		}
	}
	if v.Sub2 != nil {
		if len(v.Sub2.F01) != 3 {
			return errors.New("Sub2.F01 must be of length: 3")
		} else {
			for _, e := range v.Sub2.F01 {
				for k, e2 := range e {
					if !valid.UUID(k, 4) {
						return errors.New("Sub2.F01 must be a valid UUID")
					}
					if len(e2) < 1 {
						return errors.New("Sub2.F01 must be of length at least: 1")
					} else {
						for _, e := range e2 {
							if !valid.Zip(e, "us") {
								return errors.New("Sub2.F01 must be a valid zip code")
							}
						}
					}
				}
			}
		}
		for k, e := range v.Sub2.F02 {
			for _, e := range k {
				if !valid.IP(e, 0) {
					return errors.New("Sub2.F02 must be a valid IP")
				}
			}
			if e == nil {
				return errors.New("Sub2.F02 cannot be nil")
			}
		}
		if v.Sub2.F03 != nil {
			if v.Sub2.F03.f1 != v.Sub2.Arg1 {
				return fmt.Errorf("Sub2.F03.f1 must be equal to: %v", v.Sub2.Arg1)
			}
			if v.Sub2.F03.f2 == 0 {
				return errors.New("Sub2.F03.f2 is required")
			} else if v.Sub2.F03.f2 >= v.Arg2 {
				return fmt.Errorf("Sub2.F03.f2 must be less than: %v", v.Arg2)
			}
		}
		if !strings.HasPrefix(v.Sub2.F04, "foo") {
			return errors.New("Sub2.F04 must be prefixed with: \"foo\"")
		} else if !strings.Contains(v.Sub2.F04, "bar") {
			return errors.New("Sub2.F04 must contain substring: \"bar\"")
		} else if !strings.HasSuffix(v.Sub2.F04, "baz") && !strings.HasSuffix(v.Sub2.F04, "quux") {
			return errors.New("Sub2.F04 must be suffixed with: \"baz\" or \"quux\"")
		} else if len(v.Sub2.F04) < 9 || len(v.Sub2.F04) > 64 {
			return errors.New("Sub2.F04 must be of length between: 9 and 64 (inclusive)")
		}
		if v.Sub2.F05 != nil && *v.Sub2.F05 != "" && !valid.Email(*v.Sub2.F05) {
			return errors.New("Sub2.F05 must be a valid email address")
		}
		for _, e := range v.Sub2.F06 {
			if e != nil && !valid.Email(*e) {
				return errors.New("Sub2.F06 must be a valid email address")
			}
		}
		if v.Sub2.F07 > 0 && v.Sub2.F07 > 128 {
			return errors.New("Sub2.F07 must be less than or equal to: 128")
		}
		if v.Sub2.F08 != nil && *v.Sub2.F08 != v.Sub2.Arg1 {
			return fmt.Errorf("Sub2.F08 must be equal to: %v", v.Sub2.Arg1)
		}
		if v.Sub2.F09 >= v.Sub2.Arg2 {
			return fmt.Errorf("Sub2.F09 must be less than: %v", v.Sub2.Arg2)
		}
		if v.Sub2.F10 <= v.Sub2.Arg3 {
			return fmt.Errorf("Sub2.F10 must be greater than: %v", v.Sub2.Arg3)
		}
		if v.Sub2.Sub != nil {
			if len(v.Sub2.Sub.F01) > 10 {
				return errors.New("Sub2.Sub.F01 must be of length at most: 10")
			} else {
				for _, e := range v.Sub2.Sub.F01 {
					if !valid.Email(e) {
						return errors.New("Sub2.Sub.F01 must be a valid email address")
					}
				}
			}
			if len(v.Sub2.Sub.F02) == 0 {
				return errors.New("Sub2.Sub.F02 is required")
			} else {
				for _, e := range v.Sub2.Sub.F02 {
					if !valid.Phone(e, "us") {
						return errors.New("Sub2.Sub.F02 must be a valid phone number")
					}
				}
			}
			if v.Sub2.Sub.F03 == nil || len(*v.Sub2.Sub.F03) == 0 {
				return errors.New("Sub2.Sub.F03 is required")
			} else {
				for k := range *v.Sub2.Sub.F03 {
					if !valid.Email(k) {
						return errors.New("Sub2.Sub.F03 must be a valid email address")
					}
				}
			}
			for k, e := range v.Sub2.Sub.F04 {
				if !valid.Email(k) {
					return errors.New("Sub2.Sub.F04 must be a valid email address")
				}
				if !valid.Phone(e, "us") {
					return errors.New("Sub2.Sub.F04 must be a valid phone number")
				}
			}
			if len(v.Sub2.Sub.F05) < 5 {
				return errors.New("Sub2.Sub.F05 must be of length at least: 5")
			} else {
				for k, e := range v.Sub2.Sub.F05 {
					if !valid.Phone(k, "us") {
						return errors.New("Sub2.Sub.F05 must be a valid phone number")
					}
					if e < 21 || e > 99 {
						return errors.New("Sub2.Sub.F05 must be between: 21 and 99")
					}
				}
			}
			if v.Sub2.Sub.F06 != nil && !valid.Match(*v.Sub2.Sub.F06, `^[a-z]+\[[0-9]+\]$`) {
				return errors.New("Sub2.Sub.F06 must match the regular expression: \"^[a-z]+\\\\[[0-9]+\\\\]$\"")
			}
			if !strings.HasSuffix(v.Sub2.Sub.F07, "goodbye") {
				return errors.New("Sub2.Sub.F07 must be suffixed with: \"goodbye\"")
			}
			if v.Sub2.Sub.F08 != nil && *v.Sub2.Sub.F08 != v.Arg1 {
				return fmt.Errorf("Sub2.Sub.F08 must be equal to: %v", v.Arg1)
			}
			if v.Sub2.Sub.F09 >= v.Arg2 {
				return fmt.Errorf("Sub2.Sub.F09 must be less than: %v", v.Arg2)
			}
			if v.Sub2.Sub.F10 <= v.Arg3 {
				return fmt.Errorf("Sub2.Sub.F10 must be greater than: %v", v.Arg3)
			}
		}
	}
	return nil
}