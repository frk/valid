package as

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AS", A3: "ASM", Num: "016",
		Zip: regexp.MustCompile(`^[0-9]{5}(?:-?[0-9]{4})?$`),
	})
}
