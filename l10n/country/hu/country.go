package hu

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "HU", A3: "HUN", Num: "348",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?36)(?:20|30|70)[0-9]{7}$`),
	})
}
