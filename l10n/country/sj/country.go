package sj

import (
	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SJ", A3: "SJM", Num: "744", Zip: country.RxZip4Digits,
	})
}
