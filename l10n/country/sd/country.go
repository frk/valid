package sd

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SD", A3: "SDN", Num: "729", Zip: country.RxZip5Digits,
	})
}
