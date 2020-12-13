package rw

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "RW", A3: "RWA", Num: "646",
		Phone: regexp.MustCompile(`^(?:\+?250|0)?[7][0-9]{8}$`),
	})
}
