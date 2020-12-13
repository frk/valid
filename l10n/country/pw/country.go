package pw

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PW", A3: "PLW", Num: "585",
		Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
	})
}
