package aq

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AQ", A3: "ATA", Num: "010",
		Zip: regexp.MustCompile(`^BIQQ 1ZZ$`),
	})
}
