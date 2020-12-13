package ve

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "VE", A3: "VEN", Num: "862",
		Zip: regexp.MustCompile(`^[0-9]{4}(?:-[A-Z])?$`),
	})
}
