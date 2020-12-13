package at

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AT", A3: "AUT", Num: "040",
		Zip:   country.RxZip3Digits,
		Phone: regexp.MustCompile(`^(?:\+43|0)[0-9]{1,4}[0-9]{3,12}$`),
	})
}
