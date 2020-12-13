package me

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "ME", A3: "MNE", Num: "499",
		Zip: country.RxZip5Digits,
	})
}
