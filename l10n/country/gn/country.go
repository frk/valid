package gn

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GN", A3: "GIN", Num: "324", Zip: country.RxZip3Digits,
	})
}
