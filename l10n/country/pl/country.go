package pl

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PL", A3: "POL", Num: "616",
		Zip:   regexp.MustCompile(`^[0-9]{2}-[0-9]{3}$`),
		Phone: regexp.MustCompile(`^(?:\+?48)?[ ]?[5-8][0-9][ ]?[0-9]{3}(?:[ ]?[0-9]{2}){2}$`),
	})
}
