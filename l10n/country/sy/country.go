package sy

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SY", A3: "SYR", Num: "760",
		Phone: regexp.MustCompile(`^(?:(?:\+?963)|0)?9[0-9]{8}$`),
	})
}
