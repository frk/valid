package za

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "ZA", A3: "ZAF", Num: "710",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?27|0)[0-9]{9}$`),
	})
}
