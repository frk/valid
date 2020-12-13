package fj

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "FJ", A3: "FJI", Num: "242",
		Phone: regexp.MustCompile(`^(?:\+?679)?[ ]?[0-9]{3}[ ]?[0-9]{4}$`),
	})
}
