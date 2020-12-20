package nl

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NL", A3: "NLD", Num: "528",
		Zip:   regexp.MustCompile(`^(?i)[0-9]{4}\s?[a-z]{2}$`),
		Phone: regexp.MustCompile(`^(?:(?:(?:\+|00)?31\(0\))|(?:(?:\+|00)?31)|0)6{1}[0-9]{8}$`),
		// 'NL'+9 digits+B+2-digit company index – e.g. NL999999999B01
		VAT: regexp.MustCompile(`^NL[0-9]{9}B[0-9]{2}$`),
	})
}
