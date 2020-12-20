package at

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AT", A3: "AUT", Num: "040",
		Zip:   country.RxZip3Digits,
		Phone: regexp.MustCompile(`^(?:\+43|0)[0-9]{1,4}[0-9]{3,12}$`),
		// 'AT'+U+8 digits, – e.g. ATU99999999
		VAT: regexp.MustCompile(`^ATU[0-9]{8}$`),
	})
}
