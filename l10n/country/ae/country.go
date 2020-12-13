package ae

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AE", A3: "ARE", Num: "784",
		Phone: regexp.MustCompile(`^(?:(?:\+?971)|0)?5[024568][0-9]{7}$`),
	})
}
