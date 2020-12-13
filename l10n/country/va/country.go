package va

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "VA", A3: "VAT", Num: "336",
		Zip: regexp.MustCompile(`^00120$`),
	})
}
