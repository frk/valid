package gi

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GI", A3: "GIB", Num: "292",
		Zip: regexp.MustCompile(`^GX11 1AA$`),
	})
}
