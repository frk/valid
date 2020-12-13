package zw

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "ZW", A3: "ZWE", Num: "716",
		Phone: regexp.MustCompile(`^(?:\+263)[0-9]{9}$`),
	})
}
