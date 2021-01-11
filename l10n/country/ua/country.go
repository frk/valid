package ua

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "UA", A3: "UKR", Num: "804",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?38|8)?0[0-9]{9}$`),
		VAT:   regexp.MustCompile(`^[0-9]{12}$`),
	})
}
