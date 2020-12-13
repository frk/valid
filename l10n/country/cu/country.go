package cu

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CU", A3: "CUB", Num: "192",
		Zip: country.RxZip5Digits,
	})
}
