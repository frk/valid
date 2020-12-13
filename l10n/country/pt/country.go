package pt

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PT", A3: "PRT", Num: "620",
		Zip:   regexp.MustCompile(`^[0-9]{4}\-[0-9]{3}?$`),
		Phone: regexp.MustCompile(`^(?:\+?351)?9[1236][0-9]{7}$`),
	})
}
