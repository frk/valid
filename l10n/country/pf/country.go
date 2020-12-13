package pf

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PF", A3: "PYF", Num: "258",
		Zip: regexp.MustCompile(`^987(?:[0-8][0-9]|90)$`),
	})
}
