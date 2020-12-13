package dk

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "DK", A3: "DNK", Num: "208",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?45)?(?:[ ]?[0-9]{2}){4}$`),
	})
}
