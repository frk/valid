package ch

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CH", A3: "CHE", Num: "756",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+41|0)7[5-9][0-9]{1,7}$`),
	})
}
