package cl

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CL", A3: "CHL", Num: "152",
		Zip:   regexp.MustCompile(`^[0-9]{3}-?[0-9]{4}$`),
		Phone: regexp.MustCompile(`^(?:\+?56|0)[2-9][0-9]{8}$`),
	})
}
