package tr

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TR", A3: "TUR", Num: "792",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?90|0)?5[0-9]{9}$`),
		VAT:   regexp.MustCompile(`^[0-9]{10}$`),
	})
}
