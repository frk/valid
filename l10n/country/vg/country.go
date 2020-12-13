package vg

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "VG", A3: "VGB", Num: "092",
		Zip: regexp.MustCompile(`^VG11(?:[1-5][0-9]|60)$`),
	})
}
