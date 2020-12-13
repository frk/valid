package sh

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SH", A3: "SHN", Num: "654",
		Zip: regexp.MustCompile(`^(?:STHL|ASCN|TDCU) 1ZZ$`),
	})
}
