package cr

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CR", A3: "CRI", Num: "188",
		Zip:   regexp.MustCompile(`^[0-9]{5}(-[0-9]{4})?$`),
		Phone: regexp.MustCompile(`^(?:\+506)?[2-8][0-9]{7}$`),
	})
}
