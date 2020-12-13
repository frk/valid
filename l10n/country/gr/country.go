package gr

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GR", A3: "GRC", Num: "300",
		Zip:   regexp.MustCompile(`^[0-9]{3}[ ]?[0-9]{2}$`),
		Phone: regexp.MustCompile(`^(?:\+?30|0)?(?:69[0-9]{8})$`),
	})
}
