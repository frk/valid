package bg

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BG", A3: "BGR", Num: "100",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?359|0)?8[789][0-9]{7}$`),
		VAT:   regexp.MustCompile(`^BG[0-9]{9,10}$`),
	})
}
