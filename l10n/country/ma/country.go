package ma

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MA", A3: "MAR", Num: "504",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:(?:\+|00)212|0)[5-7][0-9]{8}$`),
	})
}
