package it

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IT", A3: "ITA", Num: "380",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?39)?[ ]?3[0-9]{2}[ ]?[0-9]{6,7}$`),
	})
}
