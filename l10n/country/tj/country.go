package tj

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TJ", A3: "TJK", Num: "762",
		Zip: country.RxZip6Digits,
	})
}
