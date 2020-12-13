package se

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SE", A3: "SWE", Num: "752",
		Zip:   regexp.MustCompile(`^[1-9][0-9]{2}\s?[0-9]{2}$`),
		Phone: regexp.MustCompile(`^(?:\+?46|0)[ \-]?7[ \-]?[02369](?:[ \-]?[0-9]){7}$`),
	})
}
