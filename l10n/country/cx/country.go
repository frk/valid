package cx

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CX", A3: "CXR", Num: "162",
		Zip: country.RxZip4Digits,
	})
}
