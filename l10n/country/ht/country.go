package ht

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "HT", A3: "HTI", Num: "332",
		Zip: regexp.MustCompile(`^HT[0-9]{4}$`),
	})
}
