package kz

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "KZ", A3: "KAZ", Num: "398",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?7|8)?7[0-9]{9}$`),
	})
}
