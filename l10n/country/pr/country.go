package pr

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PR", A3: "PRI", Num: "630",
		Zip: regexp.MustCompile(`^00[679][0-9]{2}(?:[ -][0-9]{4})?$`),
	})
}
