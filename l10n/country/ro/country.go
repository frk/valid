package ro

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "RO", A3: "ROU", Num: "642",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?4?0)[ ]?7[0-9]{2}(?:\/| |\.|\-)?[0-9]{3}(?: |\.|\-)?[0-9]{3}$`),
	})
}
