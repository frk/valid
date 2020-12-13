package us

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "US", A3: "USA", Num: "840",
		Zip: regexp.MustCompile(`^[0-9]{5}(?:-?[0-9]{4})?$`),
		Phone: regexp.MustCompile(`^(?:(?:\+?1)?[ -]?)?` +
			`(?:\([2-9][0-9]{2}\)|[2-9][0-9]{2})` +
			`[ -]?(?:[2-9][0-9]{2}[ -]?[0-9]{4})$`),
	})
}
