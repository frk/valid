package fm

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "FM", A3: "FSM", Num: "583",
		Zip: regexp.MustCompile(`^[0-9]{5}(?:-[0-9]{4})?$`),
	})
}
