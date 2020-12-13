package mh

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MH", A3: "MHL", Num: "584",
		Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
	})
}
