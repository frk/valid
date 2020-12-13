package az

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AZ", A3: "AZE", Num: "031",
		Zip:   regexp.MustCompile(`^AZ[0-9]{4}$`),
		Phone: regexp.MustCompile(`^(?:\+994|0)(?:5[015]|7[07]|99)[0-9]{7}$`),
	})
}
