package cv

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CV", A3: "CPV", Num: "132",
		Zip: country.RxZip4Digits,
	})
}
