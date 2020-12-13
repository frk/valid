package gt

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "GT", A3: "GTM", Num: "320",
		Zip: country.RxZip5Digits,
	})
}
