package ls

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LS", A3: "LSO", Num: "426",
		Zip: regexp.MustCompile(`^[0-9]{3}$`),
	})
}
