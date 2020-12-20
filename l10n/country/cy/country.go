package cy

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CY", A3: "CYP", Num: "196",
		Zip: regexp.MustCompile(`^[0-9]{4,5}$`),
		// 9 characters, last one must be a letter – e.g. CY99999999L
		VAT: regexp.MustCompile(`^CY[0-9]{8}[A-Z]$`),
	})
}
