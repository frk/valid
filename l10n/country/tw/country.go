package tw

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TW", A3: "TWN", Num: "158",
		Zip:   regexp.MustCompile(`^[0-9]{3}(?:[0-9]{2})?$`),
		Phone: regexp.MustCompile(`^(?:\+?886-?|0)?9[0-9]{8}$`),
	})
}
