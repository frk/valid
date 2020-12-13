package is

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "IS", A3: "ISL", Num: "352",
		Zip: country.RxZip3Digits,
	})
}
