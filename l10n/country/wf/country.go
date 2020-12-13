package wf

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "WF", A3: "WLF", Num: "876",
		Zip: regexp.MustCompile(`^986(?:[0-8][0-9]|90)$`),
	})
}
