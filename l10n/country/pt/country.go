package pt

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 9 digits; the last digit is the check digit. The first digit depends on
	// what the number refers to, e.g.: 1-3 are regular people, 5 are companies.
	rxvat1 := regexp.MustCompile(`^PT[0-9]{9}$`)
	rxvat2 := regexp.MustCompile(`^PT(?:(?:1|2|3|5|6|8)|(?:45|70|71|72|74|75|77|79|90|91|98|99))`)

	country.Add(country.Country{
		A2: "PT", A3: "PRT", Num: "620",
		Zip:   regexp.MustCompile(`^[0-9]{4}\-[0-9]{3}?$`),
		Phone: regexp.MustCompile(`^(?:\+?351)?9[1236][0-9]{7}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat1.MatchString(v) || !rxvat2.MatchString(v) {
				return false
			}
			v = v[2:]

			// MOD 11-2
			// https://pt.wikipedia.org/wiki/Número_de_identificação_fiscal
			sum := 0
			for i := 0; i < len(v)-1; i++ {
				num, _ := strconv.Atoi(string(v[i]))
				sum += num * (9 - i)
			}

			mod := (sum % 11)
			check, _ := strconv.Atoi(string(v[len(v)-1]))
			if mod < 2 {
				return check == 0
			}
			return check == (11 - mod)
		}),
	})
}
