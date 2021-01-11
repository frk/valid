package pe

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PE", A3: "PER", Num: "604",
		Zip:   regexp.MustCompile(`^(?:[0-9]{5})|(?:PE [0-9]{4})$`),
		Phone: regexp.MustCompile(`^(?:\+?51)?9[0-9]{8}$`),
		VAT:   regexp.MustCompile(`^[0-9]{11}$`),
	})
}
