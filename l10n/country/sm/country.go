package sm

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "SM", A3: "SMR", Num: "674",
		Zip:   regexp.MustCompile(`^4789[0-9]$`),
		Phone: regexp.MustCompile(`^(?:(?:\+378)|(?:0549)|(?:\+390549)|(?:\+3780549))?6[0-9]{5,9}$`),
		VAT:   regexp.MustCompile(`^[0-9]{5}$`),
	})
}
