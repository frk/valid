package nz

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "NZ", A3: "NZL", Num: "554",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:\+?64|0)[28][0-9]{7,9}$`),
	})
}
