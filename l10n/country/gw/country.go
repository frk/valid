package gw

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GW", A3: "GNB", Num: "624", Zip: country.RxZip4Digits,
	})
}
