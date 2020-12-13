package cc

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CC", A3: "CCK", Num: "166",
		Zip: country.RxZip4Digits,
	})
}
