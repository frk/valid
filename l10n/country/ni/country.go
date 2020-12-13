package ni

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NI", A3: "NIC", Num: "558",
		Zip: country.RxZip5Digits,
	})
}
