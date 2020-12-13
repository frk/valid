package gg

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GG", A3: "GGY", Num: "831",
		Zip:   regexp.MustCompile(`^GY[0-9]{1,2} [0-9][A-Z]{2}$`),
		Phone: regexp.MustCompile(`^(?:\+?44|0)1481[0-9]{6}$`),
	})
}
