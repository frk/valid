package jm

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "JM", A3: "JAM", Num: "388",
		Zip: regexp.MustCompile(`^[1-9]|1[0-9]|20$`),
	})
}
