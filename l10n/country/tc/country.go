package tc

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TC", A3: "TCA", Num: "796",
		Zip: regexp.MustCompile(`^TKCA 1ZZ$`),
	})
}
