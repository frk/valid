package fi

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// FI + 7 digits + check digit, e.g. FI99999999
	rxvat := regexp.MustCompile(`^FI[0-9]{8}$`)
	weigths := []int{7, 9, 10, 5, 8, 4, 2}

	country.Add(country.Country{
		A2: "FI", A3: "FIN", Num: "246",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?358|0)[ ]?(?:4(?:0|1|2|4|5|6)?|50)[ ]?(?:[0-9][ ]?){4,8}[0-9]$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			v = v[2:]

			// The check digit is calculated utilizing MOD 11-2.
			// http://tarkistusmerkit.teppovuori.fi/tarkmerk.htm#y-tunnus2
			sum := 0
			for i := 0; i < len(v)-1; i++ {
				num, _ := strconv.Atoi(string(v[i]))
				sum += num * weigths[i]
			}

			mod := (sum % 11)
			if mod == 1 {
				return false
			}

			check, _ := strconv.Atoi(string(v[len(v)-1]))
			if mod == 0 {
				return mod == check
			}
			return check == (11 - mod)
		}),
	})
}
