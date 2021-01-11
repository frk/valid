package do

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "DO", A3: "DOM", Num: "214",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?1)?8[024]9[0-9]{7}$`),
		VAT:   regexp.MustCompile(`^[0-9]{9}|[0-9]{11}$`),
	})
}
