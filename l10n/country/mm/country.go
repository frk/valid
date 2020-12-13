package mm

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "MM", A3: "MMR", Num: "104",
		Zip: country.RxZip5Digits,
	})
}
