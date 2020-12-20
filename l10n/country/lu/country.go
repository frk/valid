package lu

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LU", A3: "LUX", Num: "442",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+352)?(?:(?:6[0-9]1)[0-9]{6})$`),
		VAT:   regexp.MustCompile(`^LU[0-9]{8}$`),
	})
}
