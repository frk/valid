package gt

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GT", A3: "GTM", Num: "320",
		Zip: country.RxZip5Digits,
		// seven digits, one dash (-); one digit (like 1234567-1)
		VAT: regexp.MustCompile(`^[0-9]{7}-[0-9]$`),
	})
}
