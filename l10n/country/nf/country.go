package nf

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NF", A3: "NFK", Num: "574",
		Zip: country.RxZip4Digits,
	})
}
