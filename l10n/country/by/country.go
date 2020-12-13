package by

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BY", A3: "BLR", Num: "112",
		Zip:   regexp.MustCompile(`^2[1-4]{1}[0-9]{4}$`),
		Phone: regexp.MustCompile(`^(?:\+?375)?(?:24|25|29|33|44)[0-9]{7}$`),
	})
}
