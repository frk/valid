package sa

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SA", A3: "SAU", Num: "682",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:(?:\+?966)|0)?5[0-9]{8}$`),
	})
}
