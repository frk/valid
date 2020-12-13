package es

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "ES", A3: "ESP", Num: "724",
		Zip:   regexp.MustCompile(`^(?:5[0-2]{1}|[0-4]{1}[0-9]{1})[0-9]{3}$`),
		Phone: regexp.MustCompile(`^(?:\+?34)?[6|7][0-9]{8}$`),
	})
}
