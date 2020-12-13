package um

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "UM", A3: "UMI", Num: "581",
		Zip: regexp.MustCompile(`^96898$`),
	})
}
