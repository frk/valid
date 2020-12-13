package mf

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MF", A3: "MAF", Num: "663",
		Zip: regexp.MustCompile(`^97150$`),
	})
}
