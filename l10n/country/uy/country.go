package uy

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "UY", A3: "URY", Num: "858",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+598|0)9[1-9][0-9]{6}$`),
		VAT:   regexp.MustCompile(`^[0-9]{12}$`),
	})
}
