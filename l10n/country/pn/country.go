package pn

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PN", A3: "PCN", Num: "612",
		Zip: regexp.MustCompile(`^PCRN 1ZZ$`),
	})
}
