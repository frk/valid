package je

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "JE", A3: "JEY", Num: "832",
		Zip: regexp.MustCompile(`^JE[0-9]{1,2} [0-9][A-Z]{2}$`),
	})
}
