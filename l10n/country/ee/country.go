package ee

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "EE", A3: "EST", Num: "233",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?372)?[ ]?(?:5|8[1-4])[ ]?(?:[0-9][ ]?){6,7}$`),
		VAT:   regexp.MustCompile(`^EE[0-9]{9}$`),
	})
}
