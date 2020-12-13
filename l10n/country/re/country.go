package re

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "RE", A3: "REU", Num: "638",
		Zip:   regexp.MustCompile(`^974(?:[0-8][0-9]|90)$`),
		Phone: regexp.MustCompile(`^(?:\+?262|0|00262)[67][0-9]{8}$`),
	})
}
