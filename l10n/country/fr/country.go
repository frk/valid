package fr

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "FR", A3: "FRA", Num: "250",
		Zip:   regexp.MustCompile(`^[0-9]{2}\s?[0-9]{3}$`),
		Phone: regexp.MustCompile(`^(?:\+?33|0)[67][0-9]{8}$`),
	})
}
