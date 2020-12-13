package ec

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "EC", A3: "ECU", Num: "218",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?593|0)(?:[2-7]|9[2-9])[0-9]{7}$`),
	})
}
