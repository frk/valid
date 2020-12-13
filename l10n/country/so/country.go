package so

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SO", A3: "SOM", Num: "706",
		Zip: regexp.MustCompile(`^[A-Z]{2} [0-9]{5}$`),
	})
}
