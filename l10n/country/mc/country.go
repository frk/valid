package mc

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
	"github.com/frk/isvalid/l10n/country/fr"
)

func init() {
	country.Add(country.Country{
		A2: "MC", A3: "MCO", Num: "492",
		Zip: regexp.MustCompile(`^MC980(?:[0-9]{2})$`),
		VAT: country.StringMatcherFunc(fr.VAT),
	})
}
