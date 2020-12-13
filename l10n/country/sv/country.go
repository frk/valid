package sv

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SV", A3: "SLV", Num: "222",
		Zip: country.RxZip4Digits,
	})
}
