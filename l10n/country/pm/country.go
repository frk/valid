package pm

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PM", A3: "SPM", Num: "666",
		Zip: regexp.MustCompile(`^97500$`),
	})
}
