package pg

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PG", A3: "PNG", Num: "598",
		Zip: regexp.MustCompile(`^[0-9]{3}$`),
	})
}
