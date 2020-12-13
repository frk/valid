package im

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IM", A3: "IMN", Num: "833",
		Zip: regexp.MustCompile(`^IM[0-9]{1,2} [0-9][A-Z]{2}$`),
	})
}
