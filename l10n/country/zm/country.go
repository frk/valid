package zm

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "ZM", A3: "ZMB", Num: "894",
		Zip:   country.RxZip5Digits,
		Phone: regexp.MustCompile(`^(?:\+?26)?09[567][0-9]{7}$`),
	})
}
