package yt

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "YT", A3: "MYT", Num: "175",
		Zip: regexp.MustCompile(`^976(?:[0-8][0-9]|90)$`),
	})
}
