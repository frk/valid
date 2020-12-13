package tt

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TT", A3: "TTO", Num: "780",
		Zip: country.RxZip6Digits,
	})
}
