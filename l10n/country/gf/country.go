package gf

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GF", A3: "GUF", Num: "254",
		Zip:   regexp.MustCompile(`^973(?:[0-8][0-9]|90)$`),
		Phone: regexp.MustCompile(`^(?:\+?594|0|00594)[67][0-9]{8}$`),
	})
}
