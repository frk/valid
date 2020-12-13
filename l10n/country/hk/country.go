package hk

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "HK", A3: "HKG", Num: "344",
		Phone: regexp.MustCompile(`^(?:\+?852[\- ]?)?[456789][0-9]{3}[\- ]?[0-9]{4}$`),
	})
}
