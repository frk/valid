package pk

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PK", A3: "PAK", Num: "586",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:(?:\+92)|(?:0092))-?[0-9]{3}-?[0-9]{7}$|^[0-9]{11}$|^[0-9]{4}-[0-9]{7}$`),
	})
}
