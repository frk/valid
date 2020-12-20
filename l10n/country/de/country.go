package de

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "DE", A3: "DEU", Num: "276",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+49)?0?[1|3](?:[0|5][0-9]{2}|6(?:[23]|0[0-9]?)|7(?:[0-57-9]|6[0-9]))[0-9]{7}$`),
		// 9 digits, e.g. DE999999999
		VAT: regexp.MustCompile(`^DE[0-9]{9}$`),
	})
}
