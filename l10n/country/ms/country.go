package ms

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MS", A3: "MSR", Num: "500",
		Zip: regexp.MustCompile(`^MSR 1[1-3][0-9]{2}$`),
	})
}
