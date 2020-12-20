package dk

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 8 digits – e.g. DK99999999, last digit is check digit
	rxvat := regexp.MustCompile(`^DK[0-9]{8}$`)
	weigths := []int{2, 7, 6, 5, 4, 3, 2}

	country.Add(country.Country{
		A2: "DK", A3: "DNK", Num: "208",
		Zip:   regexp.MustCompile(`^(?:DK-)?[0-9]{4}$`),
		Phone: regexp.MustCompile(`^(?:\+?45)?(?:[ ]?[0-9]{2}){4}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			v = v[2:]

			// The check digit is calculated utilizing MOD 11-2.
			// https://web.archive.org/web/20120917151518/http://www.erhvervsstyrelsen.dk/modulus_11
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
