package country

import (
	"regexp"
	"strings"
)

// References:
// - https://wikipedia.org/wiki/ISO_3166-1_alpha-2
// - https://en.wikipedia.org/wiki/List_of_postal_codes
// - https://en.youbianku.com
// - https://en.wikipedia.org/wiki/VAT_identification_number

var ISO31661A_2 = make(map[string]Country)
var ISO31661A_3 = make(map[string]Country)

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

func Add(c Country) {
	ISO31661A_2[c.A2] = c
	ISO31661A_3[c.A3] = c
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

type StringMatcher interface {
	MatchString(v string) bool
}

type MatchStringFunc func(v string) bool

func (f MatchStringFunc) MatchString(v string) bool {
	return f(v)
}

var (
	RxZip3Digits = regexp.MustCompile(`^[0-9]{3}$`)
	RxZip4Digits = regexp.MustCompile(`^[0-9]{4}$`)
	RxZip5Digits = regexp.MustCompile(`^[0-9]{5}$`)
	RxZip6Digits = regexp.MustCompile(`^[0-9]{6}$`)
)
