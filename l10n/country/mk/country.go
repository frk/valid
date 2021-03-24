package mk

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MK", A3: "MKD", Num: "807",
		Zip: country.RxZip4Digits,
		// 15 characters, the first two positions are for the prefix
		// "MK", followed by 13 numbers – e.g. MK4032013544513
		VAT: regexp.MustCompile(`^MK[0-9]{13}$`),
	})
}
