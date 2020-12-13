package kh

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "KH", A3: "KHM", Num: "116",
		Zip: country.RxZip6Digits,
	})
}
