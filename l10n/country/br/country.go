package br

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BR", A3: "BRA", Num: "076",
		Zip:   regexp.MustCompile(`^[0-9]{5}-[0-9]{3}$`),
		Phone: regexp.MustCompile(`^(?:(?:\+?55[ ]?[1-9]{2}[ ]?)|(?:\+?55[ ]?\([1-9]{2}\)[ ]?)|(?:0[1-9]{2}[ ]?)|(?:\([1-9]{2}\)[ ]?)|(?:[1-9]{2}[ ]?))(?:(?:[0-9]{4}-?[0-9]{4})|(?:9[2-9]{1}[0-9]{3}-?[0-9]{4}))$`),
	})
}
