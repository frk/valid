package th

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "TH", A3: "THA", Num: "764",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+66|66|0)[0-9]{9}$`),
	})
}
