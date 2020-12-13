package md

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MD", A3: "MDA", Num: "498",
		Zip: regexp.MustCompile(`^MD-?[0-9]{4}$`),
	})
}
