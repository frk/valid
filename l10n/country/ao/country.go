package ao

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AO", A3: "AGO", Num: "024",
		Phone: regexp.MustCompile(`^(?:\+244)[0-9]{9}$`),
	})
}
