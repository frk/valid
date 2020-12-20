package gr

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// https://tatief.wordpress.com/2008/12/29/αλγόριθμος-του-αφμ-έλεγχος-ορθότητας/
	rxvat := regexp.MustCompile(`^(?:EL|GR)[0-9]{9}$`)
	weights := []int{256, 128, 64, 32, 16, 8, 4, 2}

	country.Add(country.Country{
		A2: "GR", A3: "GRC", Num: "300",
		Zip:   regexp.MustCompile(`^[0-9]{3}[ ]?[0-9]{2}$`),
		Phone: regexp.MustCompile(`^(?:\+?30|0)?(?:69[0-9]{8})$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}
			v = v[2:]

			sum := 0
			for i := 0; i < 8; i++ {
				n, _ := strconv.Atoi(string(v[i]))
				sum += n * weights[i]
			}

			chk, _ := strconv.Atoi(string(v[8]))
			return chk == (sum%11)%10
		}),
	})
}
