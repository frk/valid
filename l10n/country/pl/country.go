package pl

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 10 digits, the last one is a check digit; for convenience the
	// digits are separated by hyphens (xxx-xxx-xx-xx or xxx-xx-xx-xxx
	// for legal people), but formally the number consists only of digits
	rxvat := regexp.MustCompile(`^PL(?:[0-9]{10}|(?:[0-9]{3}-){2}[0-9]{2}-[0-9]{2}|[0-9]{3}-(?:[0-9]{2}-){2}[0-9]{3})$`)
	weigths := []int{6, 5, 7, 2, 3, 4, 5, 6, 7}

	country.Add(country.Country{
		A2: "PL", A3: "POL", Num: "616",
		Zip:   regexp.MustCompile(`^[0-9]{2}-[0-9]{3}$`),
		Phone: regexp.MustCompile(`^(?:\+?48)?[ ]?[5-8][0-9][ ]?[0-9]{3}(?:[ ]?[0-9]{2}){2}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}

			// remove potential hyphens
			v = strings.Map(func(r rune) rune {
				if r == '-' {
					return -1
				}
				return r
			}, v[2:])

			// https://pl.wikipedia.org/wiki/NIP
			sum := 0
			for i := 0; i < len(v)-1; i++ {
				num, _ := strconv.Atoi(string(v[i]))
				sum += num * weigths[i]
			}

			check, _ := strconv.Atoi(string(v[len(v)-1]))
			return check == (sum % 11)
		}),
	})
}
