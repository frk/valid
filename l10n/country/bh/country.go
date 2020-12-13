package bh

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "BH", A3: "BHR", Num: "048",
		Zip:   regexp.MustCompile(`^[0-9]{3,4}$`),
		Phone: regexp.MustCompile(`^(?:\+?973)?(?:3|6)[0-9]{7}$`),
	})
}
