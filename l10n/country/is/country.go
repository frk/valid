package is

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IS", A3: "ISL", Num: "352",
		Zip: country.RxZip3Digits,
		// 5 or 6 characters depending on age of the company
		VAT: regexp.MustCompile(`^[0-9]{5,6}$`),
	})
}
