package tn

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TN", A3: "TUN", Num: "788",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?216)?[2459][0-9]{7}$`),
	})
}
