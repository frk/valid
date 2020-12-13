package mx

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MX", A3: "MEX", Num: "484",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?52)?(?:1|01)?[0-9]{10,11}$`),
	})
}
