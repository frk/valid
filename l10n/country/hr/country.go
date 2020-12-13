package hr

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "HR", A3: "HRV", Num: "191",
		Zip: regexp.MustCompile(`^(?:[1-5][0-9]{4}$)`),
	})
}
