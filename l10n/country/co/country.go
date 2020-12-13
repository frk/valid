package co

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CO", A3: "COL", Num: "170",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?57)?(?:[1-8]{1}|3[0-9]{2})?[0-9]{7}$`),
	})
}
