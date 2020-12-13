package gp

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GP", A3: "GLP", Num: "312",
		Zip:   regexp.MustCompile(`^971(?:[0-8][0-9]|90)$`),
		Phone: regexp.MustCompile(`^(?:\+?590|0|00590)[67][0-9]{8}$`),
	})
}
