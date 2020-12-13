package bb

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BB", A3: "BRB", Num: "052",
		Zip: regexp.MustCompile(`^BB[0-9]{5}$`),
	})
}
