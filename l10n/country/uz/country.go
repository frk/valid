package uz

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "UZ", A3: "UZB", Num: "860",
		Zip:   country.RxZip6Digits,
		Phone: regexp.MustCompile(`^(?:\+?998)?(?:6[125-79]|7[1-69]|88|9[0-9])[0-9]{7}$`),
	})
}
