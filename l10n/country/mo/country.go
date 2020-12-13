package mo

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MO", A3: "MAC", Num: "446",
		Phone: regexp.MustCompile(`^(?:\+?853[\- ]?)?[6][0-9]{3}[\- ]?[0-9]{4}$`),
	})
}
