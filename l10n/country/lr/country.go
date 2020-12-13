package lr

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "LR", A3: "LBR", Num: "430",
		Zip: country.RxZip4Digits,
	})
}
