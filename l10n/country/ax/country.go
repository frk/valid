package ax

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AX", A3: "ALA", Num: "248",
		Zip: regexp.MustCompile(`^(?:AX-)?[0-9]{5}$`),
	})
}
