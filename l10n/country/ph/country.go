package ph

import (
	"regexp"

	"github.com/frk/isvalid/l10n/country"
)

func init() {
	country.Add(country.Country{
		A2: "PH", A3: "PHL", Num: "608",
		Zip:   country.RxZip4Digits,
		Phone: regexp.MustCompile(`^(?:09|\+639)[0-9]{9}$`),
	})
}
