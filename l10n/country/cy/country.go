package cy

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CY", A3: "CYP", Num: "196",
		Zip: regexp.MustCompile(`^[0-9]{4,5}$`),
	})
}
