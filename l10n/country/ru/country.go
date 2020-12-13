package ru

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "RU", A3: "RUS", Num: "643",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?7|8)?9[0-9]{9}$`),
	})
}
