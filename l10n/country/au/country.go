package au

import (
	"regexp"
	"strconv"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	// 11 digit number formed from a 9 digit unique identifier and two prefix
	// check digits. The two leading digits (the check digits) will be derived
	// from the subsequent 9 digits using a modulus 89 check digit calculation.
	// - https://en.wikipedia.org/wiki/VAT_identification_number
	// - https://abr.business.gov.au/Help/AbnFormat
	rxvat := regexp.MustCompile(`^[1-9][0-9]{10}$`)
	weights := []int{10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19}

	country.Add(country.Country{
		A2: "AU", A3: "AUS", Num: "036", Zip: country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?61|0)4[0-9]{8}$`),
		VAT: country.StringMatcherFunc(func(v string) bool {
			if !rxvat.MatchString(v) {
				return false
			}

			b := []byte(v)
			n, _ := strconv.Atoi(string(b[0]))
			b[0] = strconv.Itoa(n - 1)[0]

			v = string(b)

			sum := 0
			for i := 0; i < len(v); i++ {
				n, _ := strconv.Atoi(string(v[i]))
				sum += n * weights[i]
			}
			return (sum % 89) == 0
		}),
	})
}
