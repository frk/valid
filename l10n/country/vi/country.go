package vi

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "VI", A3: "VIR", Num: "850",
		Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
	})
}
