package ky

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "KY", A3: "CYM", Num: "136",
		Zip: regexp.MustCompile(`^KY[0-9]-[0-9]{4}$`),
	})
}
