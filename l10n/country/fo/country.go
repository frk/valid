package fo

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "FO", A3: "FRO", Num: "234",
		Zip:   regexp.MustCompile(`^FO-[0-9]{3}$`),
		Phone: regexp.MustCompile(`^(?:\+?298)?(?:[ ]?[0-9]{2}){3}$`),
	})
}
