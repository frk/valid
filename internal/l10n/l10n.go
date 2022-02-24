package l10n

import (
	"strings"
)

type Country struct {
	// ISO 3166-1 Alpha-2
	A2 string
	// ISO 3166-1 Alpha-3
	A3 string
	// ISO 3166-1 numeric
	Num string
	// The validator for the country's zip / postal code, will be nil
	// for countries that don't use postal codes.
	Zip StringMatcher
	// The validator for the country's phone numbers, may be nil.
	Phone StringMatcher
	//
	VAT StringMatcher
}

func Get(cc string) (c Country, ok bool) {
	if len(cc) == 2 {
		cc = strings.ToUpper(cc)
		c, ok = ISO31661A_2[cc]
		return c, ok
	}
	if len(cc) == 3 {
		cc = strings.ToUpper(cc)
		c, ok = ISO31661A_3[cc]
		return c, ok
	}
	return c, false
}

// mimics *regexp.Regexp interface
type StringMatcher interface {
	MatchString(v string) bool
}

type StringMatcherFunc func(v string) bool

func (f StringMatcherFunc) MatchString(v string) bool {
	return f(v)
}
