package hn

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "HN", A3: "HND", Num: "340",
		Zip:   regexp.MustCompile(`^(?:[0-9]{5})|(?:[A-Z]{2}[0-9]{4})$`),
		Phone: regexp.MustCompile(`^(?:\+?504)?[9|8][0-9]{7}$`),
	})
}
