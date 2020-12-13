package tm

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TM", A3: "TKM", Num: "795",
		Zip: country.RxZip6Digits,
	})
}
