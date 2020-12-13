package la

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LA", A3: "LAO", Num: "418",
		Zip: country.RxZip5Digits,
	})
}
