package bt

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BT", A3: "BTN", Num: "064",
		Zip: country.RxZip5Digits,
	})
}
