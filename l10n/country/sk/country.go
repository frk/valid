package sk

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SK", A3: "SVK", Num: "703",
		Zip:   regexp.MustCompile(`^[0-9]{3}\s?[0-9]{2}$`),
		Phone: regexp.MustCompile(`^(?:\+?421)?[ ]?[1-9][0-9]{2}(?:[ ]?[0-9]{3}){2}$`),
	})
}
