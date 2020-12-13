package fk

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "FK", A3: "FLK", Num: "238",
		Zip: regexp.MustCompile(`^FIQQ 1ZZ$`),
	})
}
