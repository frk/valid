package mv

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MV", A3: "MDV", Num: "462",
		Zip: country.RxZip5Digits,
	})
}
