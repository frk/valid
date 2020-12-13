package bn

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BN", A3: "BRN", Num: "096",
		Zip: regexp.MustCompile(`^[A-Z]{2}[0-9]{4}$`),
	})
}
