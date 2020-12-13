package rs

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "RS", A3: "SRB", Num: "688",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+3816|06)[\- 0-9]{5,9}$`),
	})
}
