package au

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AU", A3: "AUS", Num: "036", Zip: country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?61|0)4[0-9]{8}$`),
	})
}
