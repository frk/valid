package bl

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BL", A3: "BLM", Num: "652",
		Zip: regexp.MustCompile(`^97133$`),
	})
}
