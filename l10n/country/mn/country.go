package mn

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MN", A3: "MNG", Num: "496",
		Zip: country.RxZip5Digits,
	})
}
