package ug

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "UG", A3: "UGA", Num: "800",
		Phone: regexp.MustCompile(`^(?:\+?256|0)?7[0-9]{8}$`),
	})
}
