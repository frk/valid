package et

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "ET", A3: "ETH", Num: "231",
		Zip: country.RxZip4Digits,
	})
}
