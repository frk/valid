package ca

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "CA", A3: "CAN", Num: "124",
		Zip:   regexp.MustCompile(`^(?i)[ABCEGHJKLMNPRSTVXY][0-9][ABCEGHJ-NPRSTV-Z][\s\-]?[0-9][ABCEGHJ-NPRSTV-Z][0-9]$`),
		Phone: regexp.MustCompile(`^(?:(?:\+1|1)?(?: |-)?)?(?:\([2-9][0-9]{2}\)|[2-9][0-9]{2})(?: |-)?(?:[2-9][0-9]{2}(?: |-)?[0-9]{4})$`),
	})
}
