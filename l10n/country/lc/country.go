package lc

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LC", A3: "LCA", Num: "662",
		// Reference:
		// - https://stluciapostal.com/postal-codes-2/
		Zip: regexp.MustCompile(`^LC[0-9]{2}[ ]{0,2}[0-9]{3}$`),
	})
}
