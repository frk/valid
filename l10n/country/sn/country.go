package sn

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SN", A3: "SEN", Num: "686",
		Zip: country.RxZip5Digits,
	})
}
