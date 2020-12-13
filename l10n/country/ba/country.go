package ba

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BA", A3: "BIH", Num: "070",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:(?:(?:\+|00)3876)|06)(?:(?:(?:[0-3]|[5-6])[0-9]{6})|4[0-9]{7})$`),
	})
}
