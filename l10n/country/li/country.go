package li

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LI", A3: "LIE", Num: "438",
		Zip: regexp.MustCompile(`^(?:948[5-9]|949[0-7])$`),
	})
}
