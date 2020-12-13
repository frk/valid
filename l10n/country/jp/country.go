package jp

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "JP", A3: "JPN", Num: "392",
		Zip:   regexp.MustCompile(`^[0-9]{3}\-[0-9]{4}$`),
		Phone: regexp.MustCompile(`^(?:\+81[ \-]?(?:\(0\))?|0)[6789]0(?:[ \-]?[0-9]{4}){2}$`),
	})
}
