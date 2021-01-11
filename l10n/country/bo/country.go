package bo

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BO", A3: "BOL", Num: "068",
		Phone: regexp.MustCompile(`^(?:\+?591)?(?:6|7)[0-9]{7}$`),
		VAT:   regexp.MustCompile(`^[0-9]{7}$`),
	})
}
