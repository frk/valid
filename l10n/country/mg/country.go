package mg

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MG", A3: "MDG", Num: "450",
		Zip: regexp.MustCompile(`^[0-9]{3}$`),
	})
}
