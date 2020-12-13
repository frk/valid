package gs

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GS", A3: "SGS", Num: "239",
		Zip: regexp.MustCompile(`^SIQQ 1ZZ$`),
	})
}
