package ai

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "AI", A3: "AIA", Num: "660",
		Zip: regexp.MustCompile(`^AI-2640$`),
	})
}
