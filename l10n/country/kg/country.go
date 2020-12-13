package kg

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "KG", A3: "KGZ", Num: "417",
		Zip: country.RxZip6Digits,
	})
}
