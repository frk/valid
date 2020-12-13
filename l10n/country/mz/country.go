package mz

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MZ", A3: "MOZ", Num: "508",
		Zip: country.RxZip4Digits,
	})
}
