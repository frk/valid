package af

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AF", A3: "AFG", Num: "004",
		Zip: regexp.MustCompile(`^(?:[1-3][0-9]|4[0-3])(?:[0-9][1-9])$`),
	})
}
