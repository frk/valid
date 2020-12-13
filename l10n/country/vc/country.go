package vc

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "VC", A3: "VCT", Num: "670",
		Zip: regexp.MustCompile(`^VC[0-9]{4}$`),
	})
}
