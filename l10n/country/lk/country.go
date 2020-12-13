package lk

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LK", A3: "LKA", Num: "144",
		Zip: country.RxZip5Digits,
	})
}
