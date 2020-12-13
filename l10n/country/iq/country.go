package iq

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IQ", A3: "IRQ", Num: "368",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?964|0)?7[0-9]{9}$`),
	})
}
