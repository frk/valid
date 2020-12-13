package mu

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MU", A3: "MUS", Num: "480",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?230|0)?[0-9]{8}$`),
	})
}
