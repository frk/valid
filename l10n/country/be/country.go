package be

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BE", A3: "BEL", Num: "056",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?32|0)4?[0-9]{8}$`),
	})
}
