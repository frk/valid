package ir

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IR", A3: "IRN", Num: "364",
		Zip:   regexp.MustCompile(`^[0-9]{10}$`),
		Phone: regexp.MustCompile(`^(?:\+?98[\- ]?|0)9[0-39][0-9][\- ]?[0-9]{3}[\- ]?[0-9]{4}$`),
	})
}
