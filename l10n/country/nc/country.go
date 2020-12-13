package nc

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NC", A3: "NCL", Num: "540",
		Zip: regexp.MustCompile(`^988(?:[0-8][0-9]|90)$`),
	})
}
