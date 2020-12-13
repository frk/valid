package mk

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MK", A3: "MKD", Num: "807",
		Zip: country.RxZip4Digits,
	})
}
