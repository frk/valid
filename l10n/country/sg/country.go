package sg

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SG", A3: "SGP", Num: "702",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+65)?[689][0-9]{7}$`),
	})
}
