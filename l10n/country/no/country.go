package no

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NO", A3: "NOR", Num: "578",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?47)?[49][0-9]{7}$`),
	})
}
